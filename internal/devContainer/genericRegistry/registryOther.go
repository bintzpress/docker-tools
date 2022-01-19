//go:build !windows

package genericRegistry

func GetString(base string, key string, prop string) (string, error) {
	return "", nil
}
