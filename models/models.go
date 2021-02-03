package models

type ConfigModel struct {
	adbPATH        string //可执行文件的位置，eg: /path/to/adb.exe or /path/to/adb,在PATH下可留空
	networkDevices map[string]*Device
}

type Device struct {
	Name string
	Addr string
	//UUID string
}
