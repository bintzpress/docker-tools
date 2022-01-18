package commandExecutor

import (
	"github.com/bintzpress/docker-build/internal/buildConfig"
	"github.com/bintzpress/docker-build/internal/dependencyResolver"

	"errors"
	"os"
	"os/exec"
	"strings"
)

func Build(dir string, ibc buildConfig.ImageBuildConfig) error {
	var args = []string{"image",
		"build",
		ibc.Context,
		"-t",
		ibc.Tag,
	}

	for key, val := range ibc.Args {
		args = append(args, "--build-arg", key+"="+val)
	}
	dockerCmd := exec.Command("docker", args...)
	dockerCmd.Dir = dir
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr

	// Execute the command
	err := dockerCmd.Run()
	return err
}

func Pull(dir string, url string) error {
	dockerCmd := exec.Command("docker", "image", "pull", url)
	dockerCmd.Dir = dir
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr

	// Execute the command
	err := dockerCmd.Run()
	return err
}

func dockerPush(dir string, new_tag string) error {
	var args = []string{"image",
		"push",
		new_tag}

	dockerCmd := exec.Command("docker", args...)
	dockerCmd.Dir = dir
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr

	// Execute the command
	err := dockerCmd.Run()
	return err
}

func dockerRetag(dir string, original_tag string, new_tag string) error {
	var args = []string{"image",
		"tag",
		original_tag,
		new_tag}

	dockerCmd := exec.Command("docker", args...)
	dockerCmd.Dir = dir
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr

	// Execute the command
	err := dockerCmd.Run()
	return err
}

func Push(dir string, bc *buildConfig.BuildConfig, imageName string) error {
	var exists bool
	var original_tag string
	var tag_only string
	var new_tag string
	var repository string
	var regMap map[string]string
	var url string
	var index int
	var err error

	for registry, pushMap := range bc.Images[imageName].Push {
		regMap, exists = bc.Registries[registry]
		if exists {
			if bc.Images[imageName].Pull != "" {
				original_tag = bc.Images[imageName].Pull
			} else {
				original_tag = bc.Images[imageName].Build.Tag
			}

			index = strings.LastIndex(original_tag, "/")
			if index > -1 {
				// want everything after the last /
				tag_only = original_tag[:index]
			} else {
				tag_only = original_tag
			}

			url, exists = regMap["url"]
			if exists {
				new_tag = url + "/"
			} else {
				return errors.New("Missing registry url")
			}
			repository, exists = pushMap["repository"]
			if exists {
				new_tag = new_tag + repository + "/"
			}
			new_tag = new_tag + tag_only
			err = dockerRetag(dir, original_tag, new_tag)
			if err == nil {
				err = dockerPush(dir, new_tag)
			}

			if err != nil {
				break
			}
		}
	}

	return nil
}

func Execute(dir string, bc *buildConfig.BuildConfig) error {
	dependencyIterator := dependencyResolver.NewDependencyIterator()
	imageName := dependencyIterator.DependencyIteratorNext()
	var err error

	for {
		if imageName != "" {
			imageConfig, exists := bc.Images[imageName]
			if exists {
				if imageConfig.Pull != "" {
					err = Pull(dir, imageConfig.Pull)
				}
				if err == nil && imageConfig.Build.Tag != "" {
					err = Build(dir, imageConfig.Build)
				}
				if err == nil && imageConfig.Push != nil {
					err = Push(dir, bc, imageName)
				}
			}

			if err != nil {
				break
			}

			imageName = dependencyIterator.DependencyIteratorNext()
		} else {
			break
		}
	}
	return err
}
