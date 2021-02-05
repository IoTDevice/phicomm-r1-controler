package config

import (
	"bytes"
	"fmt"
	"github.com/OpenIoTHub/service-register/nettool"
	"github.com/gin-gonic/gin"
	"github.com/iotdevice/zeroconf"
	adb "github.com/mDNSService/goadb"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

type ConfigModel struct {
	ADBConfig      *adb.ServerConfig
	NetworkDevices map[string]*Device
}

type Device struct {
	Host string //Ip
	Port int    //default 5555
}

type AndroidAdbDeviceWithOpenIoTHub struct {
	*adb.Device
	Id             string
	listener       net.Listener
	zeroconfServer *zeroconf.Server
}

func (ao *AndroidAdbDeviceWithOpenIoTHub) Reg() {
	WG.Add(1)
	defer WG.Done()
	var err error
	ao.listener, err = net.Listen("tcp", ":")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(ao.listener.Addr())
	ao.RegMdns()
	ao.StartAPIServer()
}

func (ao *AndroidAdbDeviceWithOpenIoTHub) StartAPIServer() {
	herr := func(err error, c *gin.Context) {
		log.Println(err)
		c.JSON(200, gin.H{
			"code":    1,
			"message": err.Error(),
			"result":  "",
		})
	}
	var err error
	//将自己注册为OpenIoTHub设备
	r := gin.Default()
	//获取屏幕截图
	r.GET("/get-image", func(c *gin.Context) {
		_, err := ao.RunCommand("screencap", "-p", "/data/local/tmp/tmp.png")
		if err != nil {
			herr(err, c)
			return
		}
		reader, err := ao.OpenRead("/data/local/tmp/tmp.png")
		if err != nil {
			herr(err, c)
			return
		}
		bs, err := ioutil.ReadAll(reader)
		if err != nil {
			herr(err, c)
			return
		}

		extraHeaders := map[string]string{
			//"Content-Disposition":`attachment;filename="tmp.png"`,
		}
		c.DataFromReader(200, int64(len(bs)), "image/png", bytes.NewReader(bs), extraHeaders)
	})
	//执行命令
	r.GET("/do-cmd", func(c *gin.Context) {
		cmd := c.Query("cmd")
		log.Println(cmd)
		rst, err := ao.RunCommand(cmd)
		if err != nil {
			herr(err, c)
			return
		}
		c.JSON(200, gin.H{
			"code":    0,
			"message": "",
			"result":  rst,
		})
	})
	//安装apk
	r.POST("/install-apk", func(c *gin.Context) {
		file, err := c.FormFile("android.apk")
		if err != nil {
			herr(err, c)
			return
		}
		w, err := ao.Device.OpenWrite("/data/local/tmp/android.apk", 777, time.Time{})
		if err != nil {
			herr(err, c)
			return
		}
		defer w.Close()
		r, err := file.Open()
		if err != nil {
			herr(err, c)
			return
		}
		defer r.Close()
		io.Copy(w, r)
		rst, err := ao.RunCommand("/system/bin/pm", "install", "-t", "/data/local/tmp/android.apk")
		if err != nil {
			herr(err, c)
			return
		}
		c.JSON(200, gin.H{
			"code":    0,
			"message": "",
			"result":  rst,
		})
	})
	//彩灯颜色
	r.GET("/led-color", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "空接口",
		})
	})
	//按键事件，比如音量控制https://www.zhihu.com/question/50865582，其他事件https://www.cnblogs.com/zhuminghui/p/10470865.html
	r.GET("/input-keyevent", func(c *gin.Context) {
		key := c.Query("key")
		log.Println("keyevent:", key)
		rst, err := ao.RunCommand("input", "keyevent", key)
		if err != nil {
			herr(err, c)
			return
		}
		c.JSON(200, gin.H{
			"code":    0,
			"message": "",
			"result":  rst,
		})
	})
	//媒体控制
	r.GET("/media-dispatch", func(c *gin.Context) {
		key := c.Query("key")
		log.Println("keyevent:", key)
		rst, err := ao.RunCommand("media", "dispatch", key)
		if err != nil {
			herr(err, c)
			return
		}
		c.JSON(200, gin.H{
			"code":    0,
			"message": "",
			"result":  rst,
		})
	})
	err = r.RunListener(ao.listener)
	if err != nil {
		log.Println(err)
		return
	}
}

func (ao *AndroidAdbDeviceWithOpenIoTHub) RegMdns() {
	var err error
	//mdns注册
	info := nettool.MDNSServiceBaseInfo
	info["name"] = fmt.Sprintf("斐讯R1音箱-%s", ao.Id)
	info["id"] = ao.Id
	info["model"] = "com.iotserv.devices.phicomm-r1-controler"
	info["firmware-respository"] = "https://github.com/IoTDevice/phicomm-r1-controler"
	port := ao.listener.Addr().(*net.TCPAddr).Port
	log.Printf("info:%+v", info)
	log.Println("port:", port)
	ao.zeroconfServer, err = nettool.RegistermDNSService(info, port)
	if err != nil {
		log.Println(err)
		return
	}
}
