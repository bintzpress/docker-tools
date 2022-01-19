//go:build windows

package genericRegistry

func GetString(base string, key string, prop string) (string, error) {
	var value string
	var err error
	var regInfo registry.regInfo

	if base == "HKLM" {
		regInfo, err = registry.OpenKey(registry.LOCAL_MACHINE, key, registry.QUERY_VALUE)
		if err == nil {
			defer regInfo.Close()
			value, _, err = regInfo.GetStringValue(string)
		}
	}
	return value, err
}
