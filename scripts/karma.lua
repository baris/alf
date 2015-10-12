local s = require("s")

function getKarma(word)
   value = alf.brainGet("karma", word)
   if value ~= "" then
      return tonumber(value)
   end
   return nil
end

function updateWord(word, add)
   word = word:sub(0, word:len()-2)
   value = getKarma(word)
   if value then
      value = value + add
   else
      value = add
   end
   alf.brainPut("karma", word, value)
end

function updateWords(msg)
   for word in msg:gmatch("%S+") do
      if s.endswith(word, "++") then
         updateWord(word, 1)
      elseif s.endswith(word, "--") then
         updateWord(word, -1)
      end
   end
end

--------------------------

function help()
   return "karma [WORDS] -- returns the karma of the given words"
end

function processMessage ()
   msg = s.trim(alf.msg(), " ")
   updateWords(msg)

   if s.startswith(msg, "karma") then
      msg = s.trim(s.trimprefix(msg, "karma "), " ")
      karma = {}
      for word in msg:gmatch("%w+") do
         value = getKarma(word)
         if value ~= nil then
            table.insert(karma, word .. " has " .. value .. " karma point(s).")
         end
      end
      return table.concat(karma, "\n")
   end
end
