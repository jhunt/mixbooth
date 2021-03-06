FROM ubuntu:20.04
RUN apt-get update \
 && apt-get install -y icecast2 \
 && useradd radio \
 && chown -R radio:radio /etc/icecast2 /var/log/icecast2
USER radio
COPY entrypoint /usr/sbin/entrypoint
ENTRYPOINT ["entrypoint"]
CMD ["icecast2", "-c", "/etc/icecast2/icecast.xml"]
