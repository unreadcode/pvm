@echo off

echo Building...

cd "%~dp0src"

go build -o "%~dp0bin\pvm.exe" "main.go"

cd ".."

ISCC "%~dp0pvm.iss"

echo Done.

@echo on