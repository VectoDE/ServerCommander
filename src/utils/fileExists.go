package utils

import "os"

func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
