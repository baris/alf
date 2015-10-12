function processMessage ()
   if string.find(alf.msg(), "alf: attack!") then
      return alf.hubotNick .. ": image me a big f'ing weapon!\nHuh! How about this?"
   end
end
