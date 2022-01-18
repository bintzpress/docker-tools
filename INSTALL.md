# Install

You'll need to install go. This project is written in go. You can download it [here](https://go.dev/dl/).

After this follow the below directions

## For Windows

- If you just want an exe then you can run "go build ./cmd/docker-build". 
This will output docker-build.exe in current directory.

- If you want to build an installation exe then install [Inno Setup](https://jrsoftware.org/isinfo.php).
After it is installed then open powershell and change into the build\windows directory. Run ".\build.ps1".
This will create a target\windows directory in the base of the project. In there you'll find 
docker-build.exe and docker-build-setup.exe. The docker-build-setup.exe is the installer.
