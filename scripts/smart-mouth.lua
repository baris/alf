
function processMessage ()
   msg = alf.msg()
   for word in msg:gmatch("%w+") do
      value = alf.brainGet("whatis", word)
      if value ~= "" then
         msg = "Oh wait, wait! Did you say " .. word .. "?\n" .. word .. " is " .. value .. "."
         return msg
      end
   end
end
