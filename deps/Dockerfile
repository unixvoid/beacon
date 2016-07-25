FROM scratch

ENV DIFF=<DIFF>
ADD rootfs.tar.gz /
COPY config.gcfg /
COPY beacon /
COPY run.sh /

CMD ["/run.sh"]
