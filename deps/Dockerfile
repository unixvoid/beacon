FROM scratch

ADD rootfs.tar.gz /
COPY config.gcfg /
COPY beacon /
COPY run.sh /

CMD ["/run.sh"]
