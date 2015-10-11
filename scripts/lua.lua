-- for k,v in pairs(_G) do
--     str = str .. ", " .. tostring(k) .. "=" .. tostring(v) .. "\n"
-- end

function help()
   return [[lua version -- print the lua version
lua packages -- list lua packages
lua alf api -- list alf lua api]]
end

function processMessage ()
   msg = alf.msg()
   if msg == "lua version" then
      return _VERSION

   elseif msg == "lua packages" then
      pkgs = {}
      for k,v in pairs(package) do
         table.insert(pkgs, tostring(k) .. " : " .. tostring(v))
      end
      return table.concat(pkgs, "\n")

   elseif msg == "lua alf api" then
      alfapi = {}
      for k,v in pairs(alf) do
         table.insert(alfapi, tostring(k) .. " : " .. tostring(v))
      end
      return table.concat(alfapi, "\n")
   end
end
