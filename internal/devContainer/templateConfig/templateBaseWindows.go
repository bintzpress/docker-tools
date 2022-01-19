//go:build windows

package templateConfig

func GetTemplateBaseDir() (string, error) {
	var value string
	var err error
	var regInfo registry.regInfo

	regInfo, err = registry.OpenKey(registry.LOCAL_MACHINE, `Software\Bintz Press\Docker Tools`, registry.QUERY_VALUE)
	if err == nil {
		defer regInfo.Close()
		value, _, err = regInfo.GetStringValue("InstallPath")
	}
	return value, err
}
