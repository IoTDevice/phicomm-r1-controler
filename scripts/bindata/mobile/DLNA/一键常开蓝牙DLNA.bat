@echo off
title һ�������Ѷȫ��Ͱ��������      by СҶͬѧ
color 0a
echo.


echo                      һ�������Ѷȫ��Ͱ��������+DLNA
echo.
echo.
echo.                                                             by��СҶͬѧ
echo.
echo.                    ������ʹ��ǰ�������һ������������
echo.
echo.
echo                    û�з�Χ�ơ�û�з�Χ�ơ�û�з�Χ��
echo.
echo.                     ֻ�ʺϲ�������ֻ�뵱���������
echo.
echo                   ���������Ҫ�ָ�СѶ����ָ���������
echo.
echo ---------------------
set /p ip=����������IP��

adb kill-server
adb start-server

echo ��ʼͨ��ADB����

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

echo ��װDLNA

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

echo ��������˳�...
pause > nul
