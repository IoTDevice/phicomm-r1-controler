package utils

import (
	"log"
	"testing"
)

func TestExportAdb(t *testing.T) {
	err := ExportAdb("./")
	if err != nil {
		log.Println(err)
		return
	}
}

func TestExportDLNA(t *testing.T) {
	err := ExportDLNA("./")
	if err != nil {
		log.Println(err)
		return
	}
}

func TestExportSELinuxSwitch(t *testing.T) {
	err := ExportSELinuxSwitch("./")
	if err != nil {
		log.Println(err)
		return
	}
}
