FROM golang:1.13-alpine

RUN apk add --no-cache bash

ENTRYPOINT ["/entrypoint.sh"]
CMD [ "-h" ]

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY phicomm-r1-controler /bin/phicomm-r1-controler
