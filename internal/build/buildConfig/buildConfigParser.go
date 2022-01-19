package buildConfig

import (
	"github.com/bintzpress/docker-tools/internal/build/dependencyResolver"
	"github.com/bintzpress/docker-tools/internal/build/localConfig"

	"errors"
	"fmt"
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v2"
)

func LoadConfig(dir string, lc *localConfig.LocalConfig) (*BuildConfig, error) {
	dependencyResolver.Initialize()
	data, err := ioutil.ReadFile(dir + "docker-build.yml")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return nil, err
	}

	var bc BuildConfig
	err = yaml.Unmarshal(data, &bc)
	if err == nil {
		err = loadDependencies(&bc)
		if err == nil {
			err = validateConfig(&bc)
			if err == nil {
				err = replaceVars(&bc, lc)
			}
		}
	}
	return &bc, err
}

func validateConfig(bc *BuildConfig) error {
	var err error = nil

	err = dependencyResolver.Validate()
	if err == nil {
		for key, val := range bc.Images {
			if val.Pull != "" && val.Build.Tag != "" {
				err = errors.New("Cannot pull and build on same image " + key)
				break
			}

			if err == nil {
				for rkey := range bc.Images[key].Push {
					_, exists := bc.Registries[rkey]
					if !exists {
						err = errors.New("Registry referenced but not defined " + rkey)
						break
					}
				}
			}

			if err != nil {
				break
			}
		}
	}
	return err
}

func loadDependencies(bc *BuildConfig) error {
	var err error = nil
	for key, val := range bc.Images {
		if val.Depends_on == "" {
			dependencyResolver.SetNoDepends(key)
		} else {
			err = dependencyResolver.AddDependsOn(key, val.Depends_on)
			if err != nil {
				break
			}
		}
	}
	return err
}

func replaceVarsInRegistry(re *regexp.Regexp, bc *BuildConfig, lc *localConfig.LocalConfig) error {
	// replace in registry urls
	var out string
	var err error
	var registryName string
	var propName string
	var propValue string

	for registryName = range bc.Registries {
		for propName, propValue = range bc.Registries[registryName] {
			out, err = replaceVarsString(re, propValue, lc)
			if err != nil {
				break
			}
			bc.Registries[registryName][propName] = out
		}

		if err != nil {
			break
		}
	}
	return err
}

func replaceVarsInImagePush(re *regexp.Regexp, ic *ImageConfig, lc *localConfig.LocalConfig) error {
	var err error = nil

	var out string
	var registryName string
	var propName string
	var propValue string

	for registryName = range ic.Push {
		for propName, propValue = range ic.Push[registryName] {
			out, err = replaceVarsString(re, propValue, lc)
			if err != nil {
				break
			}
			ic.Push[registryName][propName] = out
		}

		if err != nil {
			break
		}
	}
	return err
}

func replaceVarsInImageBuild(re *regexp.Regexp, ic *ImageConfig, lc *localConfig.LocalConfig) error {
	var err error = nil
	var ibc ImageBuildConfig = ic.Build

	// replace in build context. easy replace since just string.
	ibc.Context, err = replaceVarsString(re, ibc.Context, lc)
	if err == nil {
		// replace in build tag. easy replace since just string.
		ibc.Tag, err = replaceVarsString(re, ibc.Tag, lc)
		if err == nil {
			// replace in build args.
			for name := range ibc.Args {
				ibc.Args[name], err = replaceVarsString(re, ibc.Args[name], lc)
				if err != nil {
					break
				}
			}
		}
	}

	// return new object to ic
	ic.Build = ibc
	return err
}

func replaceVarsInImages(re *regexp.Regexp, bc *BuildConfig, lc *localConfig.LocalConfig) error {
	var err error = nil

	var imageName string
	var imageNames []string

	var exists bool
	var imageConfig ImageConfig
	var i int

	for imageName = range bc.Images {
		imageNames = append(imageNames, imageName)
	}

	for i = 0; i < len(imageNames) && err == nil; i++ {
		imageName = imageNames[i]
		imageConfig, exists = bc.Images[imageName]
		if exists {
			// replace in pull. since pull is just a string the replace is easy.
			imageConfig.Pull, err = replaceVarsString(re, imageConfig.Pull, lc)
			if err == nil {
				// replace in pushes. this is more complex. map of map
				err = replaceVarsInImagePush(re, &imageConfig, lc)
				if err == nil {
					err = replaceVarsInImageBuild(re, &imageConfig, lc)
				}
			}

			if err == nil {
				// put new values into bc.Images
				bc.Images[imageName] = imageConfig
			}
		}
	}

	return err
}

func replaceVars(bc *BuildConfig, lc *localConfig.LocalConfig) error {
	re, err := regexp.Compile("\\$\\{[0-9a-zA-Z\\_\\-]+\\}")
	if err == nil {
		err = replaceVarsInRegistry(re, bc, lc)
		if err == nil {
			err = replaceVarsInImages(re, bc, lc)
		}
	}
	return err
}

func replaceVarsString(re *regexp.Regexp, in string, lc *localConfig.LocalConfig) (string, error) {
	matches := re.FindAllStringIndex(in, -1)

	var out string
	var err error

	err = nil
	out = ""

	if matches != nil {
		var theVar string
		var theReplace string
		var exists bool

		for _, points := range matches {
			if points[0] > 0 {
				out = string(in[:points[0]])
			} else {
				out = ""
			}
			theVar = string(in[points[0]+2 : points[1]-1])
			theReplace, exists = lc.Properties[theVar]
			if exists {
				out = out + theReplace
				if points[1] < len(in)-1 {
					out = out + string(in[points[1]:])
				}
			} else {
				err = errors.New("Missing variable")
				out = in
				break
			}
		}
	} else {
		out = in // no matches
	}

	return out, err
}
