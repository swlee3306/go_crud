@echo off & setlocal

REM chdir
cd /d %~dp0

REM Get Arg
set argc=0
set argv_msg=
set argv_ver=
set new_version=
for %%x in (%*) do Set /a argc+=1
if %argc% geq 1 (
    set argv_msg=%1
)
if %argc% geq 2 (
    set argv_ver=%2
)

REM Print Arg
echo argc: %argc%
echo argv_msg: %argv_msg%
echo argv_ver: %argv_ver%

REM Get the latest tag
git remote update
for /f "tokens=1,2,3 delims=." %%a in ('git describe --abbrev^=0 --tags') do (
    set "major=%%a"
    set "minor=%%b"
    set "patch=%%c"
)

REM Construct the new version
set /a patch+=1
set "new_version=%major%.%minor%.%patch%"

REM Revise
if [%argv_msg%]==[] (
    set "argv_msg=%date% %time%"
)
if not [%argv_ver%]==[] (
    set "new_version=%argv_ver%"
)

REM Print
echo commit...
echo New version: %new_version%
echo Commit message: %argv_msg%

REM Commit and push 22
echo pushing...
git add -A
git commit -a -m "%argv_msg%"
git tag "%new_version%"
git push origin
git push origin --tags

echo done.
timeout /nobreak 3   