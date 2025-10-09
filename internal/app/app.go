package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"

	"github.com/modern-alf/alf/internal/brain"
	"github.com/modern-alf/alf/internal/commands"
	"github.com/modern-alf/alf/internal/config"
)

// App encapsulates the modern alf runtime.
type App struct {
	cfg       config.Config
	log       *zap.SugaredLogger
	slack     *slack.Client
	socket    *socketmode.Client
	brain     brain.Store
	registry  *commands.Registry
	botUserID string
}

// New creates the runtime components necessary to run alf.
func New(cfg config.Config) (*App, error) {
	logger, err := newLogger(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	logger.Infow("booting alf", "socket_mode", cfg.SocketMode)

	store, err := brain.NewFileStore(cfg.BrainPath)
	if err != nil {
		return nil, fmt.Errorf("init brain: %w", err)
	}

	options := []slack.Option{slack.OptionDebug(strings.EqualFold(cfg.LogLevel, "debug"))}
	if cfg.SocketMode {
		options = append(options, slack.OptionAppLevelToken(cfg.AppToken))
	}

	slackClient := slack.New(cfg.BotToken, options...)

	var socketClient *socketmode.Client
	if cfg.SocketMode {
		socketClient = socketmode.New(
			slackClient,
			socketmode.OptionDebug(strings.EqualFold(cfg.LogLevel, "debug")),
		)
	}

	registry := commands.NewRegistry()
	registry.Register(commands.NewHelpCommand(registry))
	registry.Register(commands.NewPingCommand())
	registry.Register(commands.NewRememberCommand())
	registry.Register(commands.NewRecallCommand())
	registry.Register(commands.NewForgetCommand())
	registry.Register(commands.NewQuoteCommand())
	registry.Register(commands.NewAboutCommand())

	app := &App{
		cfg:      cfg,
		log:      logger,
		slack:    slackClient,
		socket:   socketClient,
		brain:    store,
		registry: registry,
	}

	if err := app.identifySelf(context.Background()); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) identifySelf(ctx context.Context) error {
	auth, err := a.slack.AuthTestContext(ctx)
	if err != nil {
		return fmt.Errorf("auth test: %w", err)
	}

	a.botUserID = auth.UserID
	a.log.Infow("connected to Slack", "team", auth.Team, "user_id", auth.UserID)
	return nil
}

// Run starts the socket mode event loop and blocks until the provided context
// is cancelled.
func (a *App) Run(ctx context.Context) error {
	if a.socket == nil {
		return errors.New("socket mode is disabled")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		err := a.socket.RunContext(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			errCh <- err
		}
		close(errCh)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			if err != nil {
				return err
			}
			return nil
		case event, ok := <-a.socket.Events:
			if !ok {
				return nil
			}
			a.handleSocketEvent(ctx, event)
		}
	}
}

func (a *App) handleSocketEvent(ctx context.Context, event socketmode.Event) {
	switch event.Type {
	case socketmode.EventTypeConnecting:
		a.log.Info("Connecting to Slack...")
	case socketmode.EventTypeConnectionError:
		a.log.Warn("Connection failed; retrying")
	case socketmode.EventTypeConnected:
		a.log.Info("Connected to Slack via Socket Mode")
	case socketmode.EventTypeEventsAPI:
		a.handleEventsAPI(ctx, event)
	case socketmode.EventTypeSlashCommand:
		a.handleSlashCommand(ctx, event)
	default:
		a.socket.Ack(*event.Request)
	}
}

func (a *App) handleEventsAPI(ctx context.Context, event socketmode.Event) {
	a.socket.Ack(*event.Request)

	eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
	if !ok {
		a.log.Warnw("failed to cast EventsAPIEvent", "data", event.Data)
		return
	}

	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		inner := eventsAPIEvent.InnerEvent
		switch ev := inner.Data.(type) {
		case *slackevents.AppMentionEvent:
			a.handleMessageEvent(ctx, ev.Channel, ev.User, ev.Text, ev.ThreadTimeStamp)
		case *slackevents.MessageEvent:
			if ev.SubType != "" || ev.BotID != "" {
				return
			}
			if ev.ChannelType == "im" {
				a.handleMessageEvent(ctx, ev.Channel, ev.User, ev.Text, ev.ThreadTimeStamp)
			}
		}
	case slackevents.AppRateLimited:
		a.log.Warn("Slack rate limited the app")
	default:
		a.log.Infow("Unhandled Events API event", "type", eventsAPIEvent.Type)
	}
}

func (a *App) handleSlashCommand(ctx context.Context, event socketmode.Event) {
	command, ok := event.Data.(slack.SlashCommand)
	if !ok {
		a.socket.Ack(*event.Request)
		return
	}

	response, found, err := a.executeCommand(ctx, command.ChannelID, command.UserID, command.Text, "")
	if err != nil {
		a.socket.Ack(*event.Request, map[string]string{"response_type": "ephemeral", "text": err.Error()})
		return
	}

	if !found {
		a.socket.Ack(*event.Request, map[string]string{"response_type": "ephemeral", "text": fmt.Sprintf("I don't know how to `%s` yet.", command.Text)})
		return
	}

	text := response.Text
	if text == "" {
		text = "👌"
	}

	payload := map[string]string{
		"response_type": "in_channel",
		"text":          text,
	}

	a.socket.Ack(*event.Request, payload)
}

func (a *App) handleMessageEvent(ctx context.Context, channel, user, text, threadTS string) {
	cleaned := a.cleanMessage(text)
	if cleaned == "" {
		return
	}

	response, found, err := a.executeCommand(ctx, channel, user, cleaned, threadTS)
	if err != nil {
		a.log.Errorw("command failed", "error", err)
		a.postMessage(channel, fmt.Sprintf("Sorry, something went wrong: %v", err), threadTS)
		return
	}

	if !found {
		a.postMessage(channel, "I'm still learning. Try `help` for what I can do.", threadTS)
		return
	}

	if response.Text != "" {
		ts := response.ThreadTS
		if ts == "" {
			ts = threadTS
		}
		a.postMessage(channel, response.Text, ts, response.UseMarkdown)
	}
}

func (a *App) executeCommand(ctx context.Context, channel, user, text, threadTS string) (commands.Response, bool, error) {
	name, args := parseCommand(text)
	if name == "" {
		return commands.Response{}, false, nil
	}

	req := commands.Request{
		Context:  ctx,
		Logger:   a.log,
		Brain:    a.brain,
		Slack:    a.slack,
		Channel:  channel,
		User:     user,
		ThreadTS: threadTS,
		Text:     text,
		Args:     args,
	}

	cmd, ok := a.registry.Find(name)
	if !ok {
		return commands.Response{}, false, nil
	}

	resp, err := cmd.Handle(req)
	if err != nil {
		return commands.Response{}, true, err
	}

	if resp.ThreadTS == "" {
		resp.ThreadTS = threadTS
	}

	return resp, true, nil
}

func (a *App) postMessage(channel, text, threadTS string, markdown ...bool) {
	options := []slack.MsgOption{slack.MsgOptionText(text, len(markdown) == 0 || markdown[0])}
	if threadTS != "" {
		options = append(options, slack.MsgOptionTS(threadTS))
	}

	if _, _, err := a.slack.PostMessage(channel, options...); err != nil {
		a.log.Errorw("failed to post message", "channel", channel, "error", err)
	}
}

func (a *App) cleanMessage(text string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return ""
	}

	if a.botUserID != "" {
		mention := fmt.Sprintf("<@%s>", a.botUserID)
		trimmed = strings.ReplaceAll(trimmed, mention, "")
	}

	return strings.TrimSpace(trimmed)
}

func parseCommand(text string) (string, []string) {
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return "", nil
	}

	name := strings.ToLower(fields[0])
	args := []string(nil)
	if len(fields) > 1 {
		args = fields[1:]
	}

	return name, args
}

func newLogger(level string) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	if err := cfg.Level.UnmarshalText([]byte(level)); err != nil {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zap.TimeEncoderOfLayout(time.RFC3339)
	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	return logger.Sugar(), nil
}
