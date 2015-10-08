-- Implements a simple string replacement function

function processMessage (message)
   if string.find(message, "cloud") then
      return string.gsub(message, "cloud", "but")
   end
end
