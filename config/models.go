package config

import (
	"bytes"
	"fmt"
	"github.com/IoTDevice/phicomm-r1-controler/utils"
	"github.com/OpenIoTHub/service-register/nettool"
	"github.com/gin-gonic/gin"
	"github.com/grandcat/zeroconf"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type AdbDeviceInfo struct {
	Serial string
}

type ConfigModel struct {
	ADBConfig      *ServerConfig
	NetworkDevices []string
}

func (cm *ConfigModel) RunAdbCmd(cmd []string) (string, error) {
	var name = "adb"
	if cm.ADBConfig.PathToAdb != "" {
		name = cm.ADBConfig.PathToAdb
	}

	cmdOut := &exec.Cmd{
		Path: name,
		Args: append([]string{name}, cmd...),
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

func (cm *ConfigModel) KillAdbServer() (string, error) {
	return cm.RunAdbCmd([]string{"kill-server"})
}

func (cm *ConfigModel) StartAdbServer() (string, error) {
	return cm.RunAdbCmd([]string{"start-server"})
}

func (cm *ConfigModel) ListDevices() (devices []*AdbDeviceInfo, err error) {
	out, err := cm.RunAdbCmd([]string{"devices"})
	if err != nil {
		log.Fatalln(err)
	}
	s := strconv.Quote(out)
	log.Println("ListDevices:")
	log.Println(s)
	out = strings.Replace(out, "List of devices attached", "", -1)
	//log.Println(out)
	var n string
	if runtime.GOOS == "windows" {
		n = "\r\n"
	} else {
		n = "\n"
	}
	out = strings.Replace(out, fmt.Sprintf("%s%s", n, n), "", -1)
	log.Println(strconv.Quote(out))
	outN := strings.Split(strings.Trim(strings.TrimSpace(out), n), n)
	log.Println(outN)
	log.Println(len(outN))
	for _, line := range outN {
		//log.Println(line)
		serialInfo := strings.SplitN(line, "\t", 2)
		//log.Println(serialInfo[0])
		if strconv.QuoteToGraphic(serialInfo[0]) == "" {
			continue
		}
		devices = append(devices, &AdbDeviceInfo{Serial: serialInfo[0]})
	}

	return
}

type AndroidAdbDeviceWithOpenIoTHub struct {
	SerialID       string
	listener       net.Listener
	zeroconfServer *zeroconf.Server
}

func (ao *AndroidAdbDeviceWithOpenIoTHub) Reg() {
	WG.Add(1)
	defer WG.Done()
	var err error
	ao.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", SingleServicePort))
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

		_, err = ao.RunAdbCommand([]string{"pull", "/data/local/tmp/tmp.png", tmpfilepath})
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
		cmdSlic := strings.Split(cmd, " ")
		rst, err := ao.RunAdbSellCommandWithSlice(cmdSlic)
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
	//执行ADB命令,比如：adb {...}
	r.GET("/do-adb-cmd", func(c *gin.Context) {
		cmd := c.Query("cmd")
		log.Println(cmd)
		cmdSlic := strings.Split(cmd, " ")
		rst, err := ao.RunAdbCommand(cmdSlic)
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
		out, err := ao.RunAdbCommand([]string{"push", tmpfilepath, r1apkfile})
		log.Println(out)

		defer func() {
			out, err := ao.RunAdbCommand([]string{"rm", r1apkfile})
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(out)
		}()

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
	//	安装DLAN

	//	安装ROOT
}

func (ao *AndroidAdbDeviceWithOpenIoTHub) RegMdns() {
	var err error
	//mdns注册
	info := nettool.GetDefaultMDNSServiceBaseInfo()
	info["name"] = fmt.Sprintf("斐讯R1音箱-%s", ao.SerialID)
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

//执行adb的命令
func (ao *AndroidAdbDeviceWithOpenIoTHub) RunAdbCommand(args []string) (string, error) {
	var name = "adb"
	if ConfigModelVar.ADBConfig.PathToAdb != "" {
		name = ConfigModelVar.ADBConfig.PathToAdb
	}

	cmdOut := &exec.Cmd{
		Path: name,
		Args: append([]string{name, "-s", ao.SerialID}, args...),
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

//执行adb shell 命令
func (ao *AndroidAdbDeviceWithOpenIoTHub) RunCommand(cmd string, args ...string) (string, error) {
	return ao.RunAdbCommand(append([]string{"shell", cmd}, args...))
}

//执行adb shell 命令
func (ao *AndroidAdbDeviceWithOpenIoTHub) RunAdbSellCommandWithSlice(args []string) (string, error) {
	return ao.RunAdbCommand(append([]string{"shell"}, args...))
}
