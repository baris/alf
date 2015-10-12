local s = require("s")

function help()
   return [[lua version -- print the lua version
lua packages -- list lua packages
lua globals -- list lua globals
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

   elseif msg == "lua globals" then
      str = ""
      for k,v in pairs(_G) do
         str = str .. ", " .. tostring(k) .. "=" .. tostring(v) .. "\n"
      end
      return str

   elseif msg == "lua alf api" then
      alfapi = {}
      for k,v in pairs(alf) do
         table.insert(alfapi, tostring(k) .. " : " .. tostring(v))
      end
      return table.concat(alfapi, "\n")

   -- Well, maybe not let this. :)
   -- elseif s.startswith(msg, "lua run") then
   --    cmd = s.trimprefix(msg, "lua run")
   --    return loadstring(cmd)()
   end
end
