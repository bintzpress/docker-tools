# Install

You'll need to install go. This project is written in go. You can download it [here](https://go.dev/dl/).

After this follow the below directions

## For Windows

- If you just want an exe then you can run "go build ./cmd/docker-builder". 
This will output docker-builder.exe in current directory.

- If you want to build an installation exe then install [Inno Setup](https://jrsoftware.org/isinfo.php).
After it is installed then open powershell and change into the build\windows directory. Run ".\build.ps1".
This will create a target\windows directory in the base of the project. In there you'll find 
docker-builder.exe and docker-builder-setup.exe. The docker-builder-setup.exe is the installer.
