@echo off
cd /d %~dp0

REM $(git describe --tags $(git rev-list --tags --max-count=1))
git describe --abbrev=0 --tags

timeout 3