#
# WARNING: Change into build\windows before running this script.
# Runs relative to the build\windows directory.
#
if (-Not (Test-Path -Path ..\..\target)) {
    New-Item -Path "..\..\" -Name "target" -ItemType "directory"
}
if (Test-Path -Path ..\..\target\windows) {
    Remove-Item -Recurse -Force ..\..\target\windows 
}
New-Item -Path "..\..\target" -Name "windows" -ItemType "directory"
Set-Location ..\..\target\windows
go build ..\..\cmd\docker-builder
ISCC.exe ..\..\build\windows\installer.iss

# Return to the directory started at
Set-Location ..\..\build\windows
