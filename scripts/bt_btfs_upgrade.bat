@echo off

echo Get BitTorrent executable path...
wmic process where name="BitTorrent.exe" get executablepath>bt_exe_path.tmp 2>nul
for /f "tokens=*" %%a in  ('more +1 "bt_exe_path.tmp"') do (
set "btpath=%%a"
)
del /f /q .\bt_exe_path.tmp 2>nul
if "%btpath%" == "" (
	echo BitTorrent executable path is not found, please make sure it is running!
    goto :end
)
:loop
if "%btpath:~-1%"==" " set "btpath=%btpath:~0,-1%"&goto loop
set btroot=%btpath:~0,-14%
:verify_btpath
echo Executable path: "%btpath%"
if not exist %btpath% (
	echo BitTorrent executable path is not found, please make sure it is runing!
	goto :end
)

echo.
echo Kill BitTorrent and BTFS task...
taskkill /f /im BitTorrent.exe
taskkill /f /im BTFS.exe

echo.
echo Backup old btfs...
set btfsroot=%btroot%btfs\
set nowtime=%time:~0,2%%time:~3,2%%time:~6,2%%time:~9,2%
set bakdir="btfs_bak_%nowtime%"
echo Backup to %bakdir%
ren %btfsroot% %bakdir%
md %btfsroot%
echo tmp > %btfsroot%btfs.exe

echo.
echo Start BitTorrent.exe...
(echo set wshell=createobject^("wscript.shell"^)
echo wshell.run"%btpath%",0,false
)>".\startbt.vbs"
.\startbt.vbs
del /f /q .\startbt.vbs 2>nul


echo.
echo Waiting for btfs upgrading...
set /a times=1
:check
for /f %%a in ('powershell -command "& {try { $response = Invoke-WebRequest http://127.0.0.1:5001/hostui/#/;$Response.StatusCode} catch {$_.Exception.Response.StatusCode.Value__}}"') do (
set statusCode=%%a
)
echo Check times: %times%
if "%statusCode%" == "200" (
	echo Success!
	goto :endcheck
)
if "%times%" == "10" (
	echo Upgrade failed, please retry!
	goto :endcheck
)
set /a times+=1
timeout /T 5 /NOBREAK
goto :check
:endcheck
echo Server started!

echo.
echo Start host ui page...
start http://127.0.0.1:5001/hostui/#/
(echo set wshell=createobject^("wscript.shell"^)
echo wscript.Sleep 1000
echo wshell.sendkeys "%^R"
echo wshell.SendKeys "%^{F5}"
)>".\refresh.vbs"
.\refresh.vbs
del /f /q .\refresh.vbs 2>nul

echo.
echo Task completed!
echo.

:end

echo Exit...
timeout /T 5