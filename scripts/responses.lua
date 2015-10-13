
local s = require("s")


function help()
   return "NICK: response QUERY == RESPONSE -- sets up response for a given string"
end

function processMessage ()
   msg = s.trim(alf.msg(), " ")

   query, response = string.match(msg, "^" .. alf.name .. ".* response (.*) == (.*)$")
   if query ~= nil and response ~= nil then
      alf.brainPut("responses", query, response)
      return "OK!"
   end

   response = alf.brainGetMatch("responses", msg)
   if response ~= nil then
      return response
   end
end
