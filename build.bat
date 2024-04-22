@echo off

echo Building...

cd "%~dp0src"

go build -o "%~dp0bin\pvm.exe" "main.go"

echo Done.

@echo on