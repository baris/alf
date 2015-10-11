
function strip(str)
   return str:match("^%s*(.-)%s*$")
end

function ltrim(str, match)
   return str:match("^"..match.."(.-)$")
end

function starts(str, match)
   if str:len() < match:len() then
      return false
   elseif str:sub(0, match:len()) == match then
      return true
   end
   return false
end

function ends(str, match)
   if str:len() < match:len() then
      return false
   elseif str:sub(str:len()-match:len()+1, str:len()) == match then
      return true
   end
   return false
end

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
      if ends(word, "++") then
         updateWord(word, 1)
      elseif ends(word, "--") then
         updateWord(word, -1)
      end
   end
end

--------------------------

function help()
   return "karma [WORDS] -- returns the karma of the given words"
end

function processMessage ()
   msg = strip(alf.msg())
   updateWords(msg)

   if starts(msg, "karma") then
      msg = strip(ltrim(msg, "karma "))
      karma = {}
      for word in msg:gmatch("%w+") do
         value = getKarma(word)
         if value ~= nil then
            table.insert(karma, word .. " has a karma of " .. value)
         end
      end
      return table.concat(karma, "\n")
   end
end
