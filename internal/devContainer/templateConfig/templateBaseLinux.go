//go:build linux

package templateConfig

func GetTemplateBaseDir() (string, error) {
	return "/usr/share/docker-tools/templates/", nil
}
