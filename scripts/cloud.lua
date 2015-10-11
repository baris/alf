-- Implements a simple string replacement function

function processMessage ()
   message = alf.msg()
   if string.find(message, "cloud") then
      return string.gsub(message, "cloud", "butt")
   end
end
