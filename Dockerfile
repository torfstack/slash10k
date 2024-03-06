FROM alpine:3.19.1

RUN mkdir /opt/scurvy10k

WORKDIR /opt/scurvy10k

COPY /bin/scurvy10k-server scurvy10k-server
COPY /assets assets

CMD ["./scurvy10k-server"]
