package services

import (
	"github.com/IoTDevice/phicomm-r1-controler/config"
	"github.com/urfave/cli/v2"
	"log"
	"os/exec"
	"path/filepath"
	"time"
)

var AndroidAdbDeviceWithOpenIoTHubMap = make(map[string]*config.AndroidAdbDeviceWithOpenIoTHub)

func Run(c *cli.Context) (err error) {
	out, err := RunAdbCommand([]string{"kill-server"})
	log.Println(out)
	if err != nil {
		log.Println(err)
	}
	//启动adb服务
	_, err = config.ConfigModelVar.StartAdbServer()
	if err != nil {
		log.Fatalln(err)
	}
	if err != nil {
		return
	}
	//连接配置文件的所有网络安卓adb设备
	if config.SingleIpPort != "" {
		ConnectOneDevice(config.SingleIpPort)
	} else {
		for _, device := range config.ConfigModelVar.NetworkDevices {
			ConnectOneDevice(device)
		}
	}
	devList, err := config.ConfigModelVar.ListDevices()
	if err != nil {
		log.Fatal(err)
	}
	for _, info := range devList {
		log.Println("List adb devices:")
		log.Printf("%+v", info)
		id := info.Serial
		log.Println("id:", id)
		AndroidAdbDeviceWithOpenIoTHubMap[id] = &config.AndroidAdbDeviceWithOpenIoTHub{
			SerialID: info.Serial,
		}
	}
	for _, do := range AndroidAdbDeviceWithOpenIoTHubMap {
		go do.Reg()
	}
	time.Sleep(time.Second * 5)
	config.WG.Wait()
	return
}

func RunAdbCommand(args []string) (string, error) {
	var name = "adb"
	if config.ConfigModelVar.ADBConfig.PathToAdb != "" {
		name = config.ConfigModelVar.ADBConfig.PathToAdb
	}

	cmdOut := &exec.Cmd{
		Path: name,
		Args: append([]string{name}, args...),
	}
	if filepath.Base(name) == name {
		if lp, err := exec.LookPath(name); err != nil {
			return "", err
		} else {
			cmdOut.Path = lp
		}
	}
	outbytes, err := cmdOut.Output()
	return string(outbytes), err
}

func ConnectOneDevice(device string) (err error) {
	log.Println("connecting :", device)
	_, err = RunAdbCommand([]string{"connect", device})
	if err != nil {
		log.Println(err)
	}
	return
}
