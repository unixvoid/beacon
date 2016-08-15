FROM scratch

ENV DIFF=<DIFF>
ADD rootfs.tar.gz /
COPY config.gcfg /
COPY beacon* /beacon
COPY run.sh /

CMD ["/run.sh"]
