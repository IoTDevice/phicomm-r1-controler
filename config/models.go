package config

import (
	"bytes"
	"fmt"
	"github.com/IoTDevice/phicomm-r1-controler/utils"
	"github.com/OpenIoTHub/service-register/nettool"
	"github.com/gin-gonic/gin"
	"github.com/iotdevice/zeroconf"
	adb "github.com/mDNSService/goadb"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type ConfigModel struct {
	ADBConfig      *adb.ServerConfig
	NetworkDevices []string
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
	r.MaxMultipartMemory = 1 << 20
	//获取屏幕截图
	r.GET("/get-image", func(c *gin.Context) {
		_, err := ao.RunCommand("screencap", "-p", "/data/local/tmp/tmp.png")
		if err != nil {
			herr(err, c)
			return
		}

		unixT := time.Now().Unix()
		tmpfilepath := path.Join(utils.GetTmpDir(), fmt.Sprintf("%d.png", unixT))

		_, err = ao.RunAdbCommand("pull", "/data/local/tmp/tmp.png", tmpfilepath)
		if err != nil {
			herr(err, c)
			return
		}
		defer func() {
			cmd := exec.Command("rm", tmpfilepath)
			out, _ := cmd.Output()
			log.Println(string(out))
		}()
		bs, err := ioutil.ReadFile(tmpfilepath)

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
		unixT := time.Now().Unix()
		tmpfilepath := path.Join(utils.GetTmpDir(), fmt.Sprintf("%d.apk", unixT))
		log.Println("tmpfilepath:", tmpfilepath)

		file, err := c.FormFile("android.apk")
		err = c.SaveUploadedFile(file, tmpfilepath)
		if err != nil {
			log.Println("err = c.SaveUploadedFile(file, tmpfilepath)")
			herr(err, c)
			return
		}
		defer func() {
			cmd := exec.Command("rm", tmpfilepath)
			out, _ := cmd.Output()
			log.Println(string(out))
		}()

		r1apkfile := fmt.Sprintf("/data/local/tmp/%d.apk", unixT)
		s, err := ao.Serial()
		var cmdname = "adb"
		if ConfigModelVar.ADBConfig.PathToAdb != "" {
			cmdname = ConfigModelVar.ADBConfig.PathToAdb
		}
		cmd := exec.Command(cmdname, "-s", s, "push", tmpfilepath, r1apkfile)
		out, err := cmd.Output()
		log.Println(string(out))
		defer ao.RunCommand("rm", r1apkfile)

		rst, err := ao.RunCommand("/system/bin/pm", "install", "-t", r1apkfile)
		if err != nil {
			log.Println("/system/bin/pm")
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
	//获取安装的软件包
	r.GET("/list-packages", func(c *gin.Context) {
		var packages []string
		rst, err := ao.RunCommand("/system/bin/pm", "list", "packages")
		if err != nil {
			herr(err, c)
			return
		}
		rst = strings.Replace(rst, "package:", "", -1)
		packages = strings.Split(rst, "\r\n")
		log.Println(packages)
		c.JSON(200, gin.H{
			"code":    0,
			"message": "",
			"result":  packages,
		})
	})
	err = r.RunListener(ao.listener)
	if err != nil {
		log.Println(err)
		return
	}
	//	查看所有软件包
	//  卸载指定软件包
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

func (ao *AndroidAdbDeviceWithOpenIoTHub) RunCommand(cmd string, args ...string) (string, error) {
	var name = "adb"
	if ConfigModelVar.ADBConfig.PathToAdb != "" {
		name = ConfigModelVar.ADBConfig.PathToAdb
	}

	cmdOut := &exec.Cmd{
		Path: name,
		Args: append([]string{name, "shell", cmd}, args...),
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

func (ao *AndroidAdbDeviceWithOpenIoTHub) RunAdbCommand(args ...string) (string, error) {
	var name = "adb"
	if ConfigModelVar.ADBConfig.PathToAdb != "" {
		name = ConfigModelVar.ADBConfig.PathToAdb
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
