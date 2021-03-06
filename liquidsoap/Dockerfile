FROM ubuntu:20.04
RUN apt-get update \
 && apt-get install -y curl ca-certificates \
 && curl -Lo /tmp/liquidsoap.deb https://github.com/savonet/liquidsoap/releases/download/v1.4.4/liquidsoap-v1.4.4_1.4.4-ubuntu-focal-amd64-1_amd64.deb \
 && apt-get install -y /tmp/liquidsoap.deb \
 && rm -f /tmp/liquidsoap.deb
USER nobody
ENTRYPOINT ["liquidsoap"]
