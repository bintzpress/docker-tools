//go:build windows

package templateConfig

import (
	"os"
	"golang.org/x/sys/windows/registry"
)

func GetTemplateBaseDir() (string, error) {
	var value string
	regInfo, err := registry.OpenKey(registry.LOCAL_MACHINE, `Software\Bintz Press\Docker Tools`, registry.QUERY_VALUE)

	if err == nil {
		defer regInfo.Close()
		value, _, err = regInfo.GetStringValue("InstallPath")
		if err == nil {
			value = value + string(os.PathSeparator) + "templates" + string(os.PathSeparator)
		} else {
			value = ""
		}
	}
	return value, err
}
