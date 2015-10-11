
function processMessage (message)
   for word in message:gmatch("%w+") do
      value = AlfBrainGet("whatis", word)
      if value ~= "" then
         msg = "Oh wait, wait! Did you say " .. word .. "? " .. word .. " is " .. value .. "."
         return msg
      end
   end
end
