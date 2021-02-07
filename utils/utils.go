package utils

import (
	"path"
	"runtime"
)

func GetTmpDir() string {
	if runtime.GOOS == "windows" {
		return path.Dir("./")
	}
	if runtime.GOOS == "android" {
		return path.Dir("/data/local/tmp/")
	}
	return path.Dir("/tmp/")
}
