# Readme

## About the Project

docker-tools is two programs
1. docker-build is a program used to pull, build, and push docker images using configuration files similar to docker-compose.yml files.
2. docker-devcontainer is a program used to create project structure for developing using a devcontainer and a stack.

## Using the Project

### docker-build

- See the test/docker-build/example/wordpress-development-images for an example of what a project looks like involving docker-build. It includes a docker-build-set.yml that defines where the docker-build.yml files are. The files tell how to pull, build, and push images.

### docker-devcontainer

- See the test/docker-devcontainer/example/wordpress-development-stack for an example of a stack. You need to have a stack you want to use with your devcontainer. 

- You can do "docker-devcontainer list" to see what templates are in your app directory. They are in the templates directory. Currently I just have wordpress-developer.

- You create a directory for your project and move into it.

- An example of running the docker-devcontainer is below.

docker-devcontainer init --template wordpress-developer  --name newproject --title "New Project" --stack ..\..\stacks\wordpress-development-stack --author "Brian Bintz"

- After you have initialized your project you can open the stack folder in VS Code. It will ask to open as a devcontainer but first you need to rename the sample .env files in the .devcontainer directory. You may need to edit them. These files should not be checked into git as they are specific for your environment.

- Once you have setup your .env files you can click on the bottom left in VS Code and select "Reopen in container". This will open the stack in a devcontainer. For the wordpress stack you will be in the projects directory on the wordpress-developer container.

## About Me

My name is Brian Bintz. I am a freelance writer, developer, and trainer. Check out my profile on [GitHub](https://github.com/bintzpress) for more details on me. If you have any issues with this project please contact me at brian@bintzpress.com.
