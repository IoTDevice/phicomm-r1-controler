package services

import (
	"fmt"
	"github.com/IoTDevice/phicomm-r1-controler/config"
	adb "github.com/mDNSService/goadb"
	"github.com/urfave/cli/v2"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
	adbClient, err := adb.NewWithConfig(*config.ConfigModelVar.ADBConfig)
	if err != nil {
		return
	}
	//连接配置文件的所有网络安卓adb设备
	for _, device := range config.ConfigModelVar.NetworkDevices {
		var ip string
		var port int
		if sn := strings.SplitN(device, ":", 2); strings.Contains(device, ":") && len(sn) == 2 {
			ip = sn[0]
			port, err = strconv.Atoi(sn[1])

		} else {
			ip = device
			port = 5555
		}
		err := adbClient.Connect(ip, port)
		if err != nil {
			log.Println(err)
		}
	}
	devList, err := adbClient.ListDevices()
	if err != nil {
		log.Fatal(err)
	}
	for _, info := range devList {
		log.Println("List adb devices:")
		log.Printf("%+v", info)
		dev := adbClient.Device(adb.DeviceWithSerial(info.Serial))
		//rst, err := dev.RunCommand("ls")
		//if err != nil {
		//	log.Println(err)
		//	continue
		//}
		//log.Println(rst)
		id := fmt.Sprintf("%s-%s", info.Model, info.Serial)
		log.Println("id:", id)
		AndroidAdbDeviceWithOpenIoTHubMap[id] = &config.AndroidAdbDeviceWithOpenIoTHub{
			Device:   dev,
			SerialID: info.Serial,
		}
		log.Println(id, dev.String())
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
