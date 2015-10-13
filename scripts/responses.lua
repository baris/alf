
local s = require("s")


function help()
   return "NICK: response QUERY == RESPONSE -- sets up response for a given string"
end

function processMessage ()
   msg = s.trim(alf.msg(), " ")

   response = alf.brainGet("responses", msg)
   if response ~= "" then
      return response
   end

   query, response = string.match(msg, "^" .. alf.name .. ".* response (.*) == (.*)$")
   if query == nil or response == nil then
      return ""
   end
   alf.brainPut("responses", query, response)
   return "OK!"
end
