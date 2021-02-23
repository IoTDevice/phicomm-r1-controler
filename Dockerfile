FROM ubuntu:latest

RUN apt-get update
RUN apt-get install software-properties-common -y
RUN add-apt-repository ppa:nilarimogard/webupd8
RUN apt-get update
RUN apt-get install android-tools-adb -y

ENTRYPOINT ["/entrypoint.sh"]
CMD [ "-h" ]

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY phicomm-r1-controler /bin/phicomm-r1-controler
