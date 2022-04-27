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
set hms=%time:~0,2%%time:~3,2%%time:~6,2%%time:~9,2%
ren %btfsroot% btfs_bak_%hms%
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
timeout /T 19 /NOBREAK

echo.
echo Start http://127.0.0.1:5001/hostui/#/...
start http://127.0.0.1:5001/hostui/#/
(echo set wshell=createobject^("wscript.shell"^)
echo wscript.Sleep 2000
echo wshell.sendkeys "%^R"
echo wshell.SendKeys "%^{F5}"
)>".\refresh.vbs"
.\refresh.vbs
del /f /q .\refresh.vbs 2>nul

echo.
echo Task completed!
echo.

:end
pause