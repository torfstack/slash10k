

FROM alpine:3.21.0

RUN mkdir /opt/slash10k

WORKDIR /opt/slash10k

COPY /bin/slash10k slash10k

CMD ["./slash10k"]
