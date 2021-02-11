@echo off
title 一键禁用斐讯全家桶常开蓝牙      by 小叶同学
color 0a
echo.


echo                      一键禁用斐讯全家桶常开蓝牙+DLNA
echo.
echo.
echo.                                                             by：小叶同学
echo.
echo.                    ！！！使用前请先配对一次蓝牙！！！
echo.
echo.
echo                    没有氛围灯、没有氛围灯、没有氛围灯
echo.
echo.                     只适合不想折腾只想当蓝牙音响的
echo.
echo                   如果后续需要恢复小讯，请恢复出厂设置
echo.
echo ---------------------
set /p ip=请输入音箱IP：

adb kill-server
adb start-server

echo 开始通过ADB连接

adb connect %ip%


adb shell /system/bin/pm  uninstall --user 0 com.phicomm.speaker.airskill
adb shell /system/bin/pm  uninstall --user 0 com.phicomm.speaker.player
adb shell /system/bin/pm  uninstall --user 0 com.phicomm.speaker.exceptionreporter
adb shell /system/bin/pm  uninstall --user 0 com.phicomm.speaker.ijetty
adb shell /system/bin/pm  uninstall --user 0 com.android.keychain
adb shell /system/bin/pm  uninstall --user 0 com.phicomm.speaker.netctl
adb shell /system/bin/pm  uninstall --user 0 com.phicomm.speaker.otaservice
adb shell /system/bin/pm  uninstall --user 0 com.phicomm.speaker.systemtool
adb shell /system/bin/pm  uninstall --user 0 com.phicomm.speaker.device
adb shell /system/bin/pm  uninstall --user 0 com.android.providers.downloads
adb shell /system/bin/pm  uninstall --user 0 com.android.location.fused
adb shell /system/bin/pm  uninstall --user 0 com.android.inputdevices
adb shell /system/bin/pm  uninstall --user 0 com.android.server.telecom
adb shell /system/bin/pm  uninstall --user 0 com.android.providers.telephony
adb shell /system/bin/pm  uninstall --user 0 com.android.vpndialogs
adb shell /system/bin/pm  uninstall --user 0 com.phicomm.speaker.productiontest
adb shell /system/bin/pm  uninstall --user 0 com.phicomm.speaker.bugreport
adb shell /system/bin/pm uninstall com.droidlogic.mediacenter

adb shell settings put secure install_non_market_apps 1

adb push dlna.apk /mnt/internal_sd/

echo 安装DLNA

adb shell /system/bin/pm install -r /mnt/internal_sd/dlna.apk

adb shell rm /mnt/internal_sd/dlna.apk

adb shell am start -n com.droidlogic.mediacenter/.MediaCenterActivity

adb shell input tap 100 150
adb shell input tap 100 150
adb shell input tap 500 100
adb shell input tap 500 150

adb shell input tap 100 200
adb shell input tap 500 100
adb shell input tap 500 150
adb reboot

:end

echo 按任意键退出...
pause > nul
