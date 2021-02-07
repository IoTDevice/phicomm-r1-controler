package config

import (
	"fmt"
	adb "github.com/mDNSService/goadb"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var WG sync.WaitGroup

var ConfigFileName = "phicomm-r1-controler.yaml"
var ConfigFilePath = fmt.Sprintf("./%s", ConfigFileName)
var ConfigModelVar = &ConfigModel{
	ADBConfig:      &adb.ServerConfig{PathToAdb: ""},
	NetworkDevices: []string{"192.168.123.146:5555"},
}

//将配置写入指定的路径的文件
func WriteConfigFile(ConfigMode *ConfigModel, path string) (err error) {
	configByte, err := yaml.Marshal(ConfigMode)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	if ioutil.WriteFile(path, configByte, 0644) == nil {
		return
	}
	return
}

func InitConfigFile() {
	//	生成配置文件模板
	err := os.MkdirAll(filepath.Dir(ConfigFilePath), 0644)
	if err != nil {
		return
	}
	err = WriteConfigFile(ConfigModelVar, ConfigFilePath)
	if err == nil {
		fmt.Println("config created")
		return
	}
	log.Fatalln("写入配置文件模板出错，请检查本程序是否具有写入权限！或者手动创建配置文件。")
}

func UseConfigFile() {
	//配置文件存在
	log.Println("使用的配置文件位置：", ConfigFilePath)
	content, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	err = yaml.Unmarshal(content, ConfigModelVar)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	return
}

func LoadSnapcraftConfigPath() {
	//是否是snapcraft应用，如果是则从snapcraft指定的工作目录保存配置文件
	appDataPath, havaAppDataPath := os.LookupEnv("SNAP_USER_DATA")
	if havaAppDataPath {
		ConfigFilePath = filepath.Join(appDataPath, ConfigFileName)
	}
}
