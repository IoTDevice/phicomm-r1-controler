#!/usr/bin/env bash

echo postinstall.sh
systemctl enable phicomm-r1-controler
systemctl start phicomm-r1-controler