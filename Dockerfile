FROM golang:1.17-bullseye
COPY . /opt/short_link_server
WORKDIR /opt/short_link_server/build
RUN apt update \
    && apt install -y sqlite3 libsqlite3-dev librocksdb-dev libz-dev libbz2-dev libsnappy-dev liblz4-dev libzstd-dev \
    && cd /opt/short_link_server \
    && ./build.sh
EXPOSE 8080
ENTRYPOINT ["./server", "--config=config.yaml"]
VOLUME /opt/short_link_server/build/rocksdb
VOLUME /opt/short_link_server/build/sqlite.db