package services

import (
	"fmt"
	"github.com/IoTDevice/phicomm-r1-controler/config"
	adb "github.com/mDNSService/goadb"
	"github.com/urfave/cli/v2"
	"log"
	"time"
)

var AndroidAdbDeviceWithOpenIoTHubMap = make(map[string]*config.AndroidAdbDeviceWithOpenIoTHub)

func Run(c *cli.Context) (err error) {
	//启动adb服务
	adbClient, err := adb.NewWithConfig(*config.ConfigModelVar.ADBConfig)
	if err != nil {
		return
	}
	//连接配置文件的所有网络安卓adb设备
	for _, device := range config.ConfigModelVar.NetworkDevices {
		err := adbClient.Connect(device.Host, device.Port)
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
		AndroidAdbDeviceWithOpenIoTHubMap[fmt.Sprintf("%s-%s", info.Model, info.Serial)] = &config.AndroidAdbDeviceWithOpenIoTHub{
			Device: dev,
			Id:     id,
		}
	}
	for _, do := range AndroidAdbDeviceWithOpenIoTHubMap {
		go do.Reg()
	}
	time.Sleep(time.Second * 5)
	config.WG.Wait()
	return
}
