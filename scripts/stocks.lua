local s = require("s")

function help()
   return "NICK: stock SYMBOL -- fetch latest stock value"
end

function processMessage ()
   msg = s.trim(alf.msg(), " ")
   symbol, response = string.match(msg, "^" .. alf.name .. ".*stock (.*)$")
   url = "http://download.finance.yahoo.com/d/quotes.csv?f=nsl1op&s="..symbol
   response = http.GET(url)
   return response
end
