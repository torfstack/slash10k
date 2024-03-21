FROM alpine:3.19.1

RUN mkdir /opt/slash10k

WORKDIR /opt/slash10k

COPY /bin/slash10k-server slash10k-server
COPY /assets assets
COPY /sql/migrations sql/migrations

CMD ["./slash10k-server"]
