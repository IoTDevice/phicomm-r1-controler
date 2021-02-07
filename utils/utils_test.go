package utils

import (
	"log"
	"testing"
)

func TestGetTmpDir(t *testing.T) {
	log.Println(GetTmpDir())
}
