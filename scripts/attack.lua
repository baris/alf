function help()
   return "NICK: attack! -- you're playing with fire!"
end

function processMessage ()
   if string.find(alf.msg(), alf.name .. ": attack!") then
      return alf.hubotNick .. ": image me a big f'ing weapon!\nHuh! How about this?"
   end
end
