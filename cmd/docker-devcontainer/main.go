package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bintzpress/docker-tools/internal/copy"
	"github.com/bintzpress/docker-tools/internal/devContainer/templateConfig"
)

type CommandConfig struct {
	Name        string
	Title       string
	Stack       string
	Action      string
	Template    string
	Destination string
	Author      string
}

func parseArguments(config *CommandConfig) error {
	var err error = nil
	var templateName string
	var templateDir string
	var out string

	if len(os.Args) < 2 {
		err = errors.New("Must provide action to perform")
	} else {
		config.Action = os.Args[1]
		if config.Action != "init" && config.Action != "list" {
			err = errors.New("Only actions supported are init and list")
		} else {
			var i int
			for i = 2; i < len(os.Args); i = i + 2 {
				switch os.Args[i] {
				case "--name":
					if len(os.Args) < i+1 {
						err = errors.New("Missing value for flag --name")
					} else {
						config.Name = os.Args[i+1]
					}
					break
				case "--title":
					if len(os.Args) < i+1 {
						err = errors.New("Missing value for flag --title")
					} else {
						config.Title = os.Args[i+1]
					}
					break
				case "--stack":
					if len(os.Args) < i+1 {
						err = errors.New("Missing value for flag --stack")
					} else {
						config.Stack = os.Args[i+1]
					}
					break
				case "--template":
					if len(os.Args) < i+1 {
						err = errors.New("Missing value for flag --template")
					} else {
						templateName = os.Args[i+1]
					}
					break
				case "--templateDir":
					if len(os.Args) < i+1 {
						err = errors.New("Missing value for flag --templateDir")
					} else {
						templateDir = os.Args[i+1]
					}
					break
				case "--destination":
					if len(os.Args) < i+1 {
						err = errors.New("Missing value for flag --destination")
					} else {
						config.Destination = os.Args[i+1]
					}
					break
				case "--author":
					if len(os.Args) < i+1 {
						err = errors.New("Missing value for flag --author")
					} else {
						config.Author = os.Args[i+1]
					}
					break
				default:
					err = errors.New("Invalid arg " + os.Args[i])
					break
				}

				if err != nil {
					break
				}
			}

			if config.Action == "init" {
				if config.Name == "" || config.Action == "" || config.Stack == "" || config.Author == "" {
					err = errors.New("Missing required arg. Must provide action, --name, --stack, and --author")
				}

				if err == nil {
					if templateDir != "" {
						out, err = verifyCustomTemplate(templateDir)
						if err == nil {
							config.Template = out
						}
					} else if templateName != "" {
						out, err = verifyIncludedTemplate(os.Args[0], templateName)
						if err == nil {
							config.Template = out
						}
					} else {
						err = errors.New("Either --template or --templateDir is required")
					}
				}

				if err == nil {
					if config.Title == "" {
						config.Title = config.Name
					}
					if config.Destination == "" {
						config.Destination = "." + string(os.PathSeparator)
					}
					if strings.Contains(config.Stack, string(os.PathSeparator)) {
						out, err = filepath.Abs(config.Stack)
						if err == nil {
							config.Stack = ensureEndsWithPathSeparator(out)
						} else { // we will just use what we have
							config.Stack = ensureEndsWithPathSeparator(config.Stack)
							err = nil // ignore the error
						}
					}
					config.Template = ensureEndsWithPathSeparator(config.Template)
					config.Destination = ensureEndsWithPathSeparator(config.Destination)
				}
			}
		}
	}
	return err
}

func listTemplates() error {
	var fds []fs.FileInfo
	dir, err := templateConfig.GetTemplateBaseDir()
	if err == nil {
		if fds, err = ioutil.ReadDir(dir); err != nil {
			return err
		}

		for _, fd := range fds {
			if fd.IsDir() {
				fmt.Println(fd.Name())
			}
		}
	}
	return err
}

func verifyIncludedTemplate(exefp string, tn string) (string, error) {
	var err error
	var baseDir string
	var out string

	baseDir, err = templateConfig.GetTemplateBaseDir()
	if err == nil {
		_, err = os.Stat(baseDir + tn)
		if err == nil {
			out = baseDir + tn
		}
	}

	return out, err
}

func verifyCustomTemplate(filefp string) (string, error) {
	_, err := os.Stat(filefp)
	if err == nil {
		return filefp, err
	} else {
		return "", err
	}
}

func ensureEndsWithPathSeparator(in string) string {
	if !strings.HasSuffix(in, string(os.PathSeparator)) {
		return in + string(os.PathSeparator)
	} else {
		return in
	}
}

func copyRoot(config *CommandConfig) error {
	var err error = nil

	err = copy.DirCopy(config.Template+string(os.PathSeparator)+"root", config.Destination)
	return err
}

func replaceVars(in string, re *regexp.Regexp, name string) (string, error) {
	matches := re.FindAllStringIndex(in, -1)

	var out string
	var err error

	err = nil
	out = ""

	if matches != nil {
		for _, points := range matches {
			if points[0] > 0 {
				out = string(in[:points[0]])
			} else {
				out = ""
			}
			out = out + name
			if points[1] < len(in)-1 {
				out = out + string(in[points[1]:])
			}
		}
	} else {
		out = in // no matches
	}

	return out, err
}

func renameFiles(cc *CommandConfig, tc *templateConfig.TemplateConfig) error {
	var err error = nil
	var src string
	var replacement string
	var re *regexp.Regexp
	var dst string

	re, err = regexp.Compile(`\{\%= name \%\}`)
	if err == nil {
		for src, replacement = range tc.Rename {
			dst, err = replaceVars(replacement, re, cc.Name)
			if err == nil {
				err = os.Rename(cc.Destination+string(os.PathSeparator)+src,
					cc.Destination+string(os.PathSeparator)+dst)
			}

			if err != nil {
				break
			}
		}
	}
	return err
}

func replaceVarsInText(re *regexp.Regexp, in string, replacements *map[string]string) (string, error) {
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
			theVar = string(in[points[0]+4 : points[1]-3])
			theReplace, exists = (*replacements)[theVar]
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

func replaceVarsInFile(filefp string, replacements *map[string]string) error {
	var err error
	var re *regexp.Regexp
	var newText string
	var inBuffer []byte
	var outBuffer []byte
	var newBytes []byte

	fmt.Println("Replacing variables in " + filefp)
	re, err = regexp.Compile(`\{\%= [0-9a-zA-Z\-\_]+ \%\}`)
	if err == nil {
		inBuffer, err = ioutil.ReadFile(filefp)
		if err == nil {
			textScanner := bufio.NewScanner(strings.NewReader(string(inBuffer)))
			textScanner.Split(bufio.ScanLines)
			for textScanner.Scan() {
				newText, err = replaceVarsInText(re, textScanner.Text(), replacements)
				if err == nil {
					newBytes = []byte(newText + "\n")
					outBuffer = append(outBuffer, newBytes...)
				} else {
					break
				}
			}

			if err == nil {
				os.WriteFile(filefp, outBuffer, 0644)
			}
		}
	}
	return err
}

func RecursiveReplaceTextInDir(dir string, action func(string, *map[string]string) error, replacements *map[string]string) error {
	var err error
	var fds []os.FileInfo

	if fds, err = ioutil.ReadDir(dir); err != nil {
		return err
	}

	for _, fd := range fds {
		dirfp := dir + string(os.PathSeparator) + fd.Name()

		if fd.IsDir() {
			if err = RecursiveReplaceTextInDir(dirfp, action, replacements); err != nil {
				fmt.Println(err) // continues if can't walk part
			}
		} else {
			err = action(dirfp, replacements)
		}

		if err != nil {
			break
		}
	}

	return err
}

func makeDockerComposeFilesReplacement(stackDir string) (string, error) {
	var err error
	var fds []os.FileInfo
	var re *regexp.Regexp
	var dir string
	var pathWithLocalEnv string

	if strings.Contains(stackDir, string(os.PathSeparator)) {
		// we have a full path
		dir = stackDir
		pathWithLocalEnv = dir
	} else {
		// we have a name of a stack
		dir = os.Getenv("DockerToolsStackPath") + string(os.PathSeparator) + stackDir + string(os.PathSeparator)
		pathWithLocalEnv = `${localEnv:DockerToolsStackPath}` + string(os.PathSeparator) + stackDir + string(os.PathSeparator)
	}
	if fds, err = ioutil.ReadDir(dir); err != nil {
		return "", err
	}
	var out string

	re, err = regexp.Compile(`docker-compose[a-zA-Z0-9\_\-]+\.yml`)
	for _, fd := range fds {
		if re.MatchString(fd.Name()) {
			if out == "" {
				out = `"` + pathWithLocalEnv + fd.Name() + `"`
			} else {
				out = out + ",\n\t\t\"" + pathWithLocalEnv + fd.Name() + "\""
			}
		}
	}

	if out != "" {
		out = out + ","
	}
	return out, err
}

func copyEnvFiles(destDir string, srcDir string) error {
	var err error
	var fds []os.FileInfo
	var re *regexp.Regexp
	var dirfp string

	if !strings.Contains(srcDir, string(os.PathSeparator)) {
		// we have a name of a stack
		srcDir = os.Getenv("DockerToolsStackPath") + string(os.PathSeparator) + srcDir + string(os.PathSeparator)
	}

	if fds, err = ioutil.ReadDir(srcDir); err != nil {
		return err
	}

	re, err = regexp.Compile(`.env(?:[\w\_\-]+)*`)
	for _, fd := range fds {
		dirfp = srcDir + fd.Name()

		if re.MatchString(fd.Name()) {
			fmt.Println("Found env file " + fd.Name())
			err = copy.FileCopy(dirfp, destDir+"stack"+string(os.PathSeparator)+fd.Name())
		}
	}

	return err
}

func main() {
	var err error
	var tc *templateConfig.TemplateConfig
	var replacements map[string]string
	var out string

	cc := new(CommandConfig)
	err = parseArguments(cc)
	if err == nil {
		if cc.Action == "list" {
			err = listTemplates()
		} else if cc.Action == "init" {
			replacements = make(map[string]string)
			replacements["name"] = cc.Name
			replacements["title"] = cc.Title
			replacements["author"] = cc.Author

			out, err = makeDockerComposeFilesReplacement(cc.Stack)
			if err == nil {
				replacements["docker-compose-files"] = out
				tc, err = templateConfig.LoadConfig(cc.Template)
				if err == nil || errors.Is(err, os.ErrNotExist) {
					err = copyRoot(cc)
					if err == nil {
						if tc != nil {
							err = renameFiles(cc, tc)
						}
						if err == nil {
							err = RecursiveReplaceTextInDir(cc.Destination, replaceVarsInFile, &replacements)
							if err == nil {
								err = copyEnvFiles(cc.Destination, cc.Stack)
							}
						}
					}
				}
			}
		}
	}

	if err != nil {
		fmt.Println(err)
	}
}
