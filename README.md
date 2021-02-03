# phicomm-r1-controler
[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-white.svg)](https://snapcraft.io/phicomm-r1-controler)

```sh
phicomm-r1-controler -c /path/to/config/file/phicomm-r1-controler.yaml
```
or just:
```
phicomm-r1-controler
```
(use default config file: ./phicomm-r1-controler.yaml)

Here are the steps for each of them:

## Install the pre-compiled binary

**homebrew tap** :

```sh
$ brew install OpenIoTHub/tap/phicomm-r1-controler
```

**homebrew** (may not be the latest version):

```sh
$ brew install phicomm-r1-controler
```

**snapcraft**:

```sh
$ sudo snap install phicomm-r1-controler
```
config file path: /root/snap/phicomm-r1-controler/current/phicomm-r1-controler.yaml

edit config file then:
```sh
sudo snap restart phicomm-r1-controler
```

**scoop**:

```sh
$ scoop bucket add OpenIoTHub https://github.com/OpenIoTHub/scoop-bucket.git
$ scoop install phicomm-r1-controler
```

**deb/rpm**:

Download the `.deb` or `.rpm` from the [releases page][releases] and
install with `dpkg -i` and `rpm -i` respectively.

config file path: /etc/phicomm-r1-controler/phicomm-r1-controler.yaml

edit config file then:
```sh
sudo systemctl restart phicomm-r1-controler
```

**manually**:

Download the pre-compiled binaries from the [releases page][releases] and
copy to the desired location.

[releases]: https://github.com/IoTDevice/phicomm-r1-controler/releases

