package utils

import (
	"github.com/IoTDevice/phicomm-r1-controler/assets"
	"io/ioutil"
	"path/filepath"
)

func ExportFiles(frompath, topath string) (err error) {
	data, err := assets.Asset(frompath)
	if err != nil {
		return
	}
	return ioutil.WriteFile(topath, data, 0777)
}

// scripts/bindata/host/adb/AdbWinApi.dll
// scripts/bindata/host/adb/AdbWinUsbApi.dll
// scripts/bindata/host/adb/adb.exe
func ExportAdb(topath string) (err error) {
	filename := "adb.exe"
	frompath := "scripts/bindata/host/adb/"
	err = ExportFiles(filepath.Join(frompath, filename), filepath.Join(topath, filename))
	if err != nil {
		return
	}
	filename = "AdbWinApi.dll"
	err = ExportFiles(filepath.Join(frompath, filename), filepath.Join(topath, filename))
	if err != nil {
		return
	}
	filename = "AdbWinUsbApi.dll"
	err = ExportFiles(filepath.Join(frompath, filename), filepath.Join(topath, filename))
	if err != nil {
		return
	}
	return
}

// scripts/bindata/mobile/DLNA/dlna.apk
// scripts/bindata/mobile/DLNA/一键常开蓝牙DLNA.bat
func ExportDLNA(topath string) (err error) {
	filename := "dlna.apk"
	frompath := "scripts/bindata/mobile/DLNA/"
	err = ExportFiles(filepath.Join(frompath, filename), filepath.Join(topath, filename))
	if err != nil {
		return
	}
	return
}

// scripts/bindata/mobile/ROOT/TheSELinuxSwitch7.0.0.apk
func ExportSELinuxSwitch(topath string) (err error) {
	filename := "TheSELinuxSwitch7.0.0.apk"
	frompath := "scripts/bindata/mobile/ROOT/"
	err = ExportFiles(filepath.Join(frompath, filename), filepath.Join(topath, filename))
	if err != nil {
		return
	}
	return
}

// scripts/bindata/scrcpy/scrcpy-server
