# ShortLinkServer

ShortLinkServer is a http server which can convert url to short url. e.g. ```https://www.bilibili.com/video/BV1Sx411T7QQ``` -> ```localhost:8080/ip3Thg```. 
server is written by [Go][1] and use [gin][2] framework, [RocksDB][3] as kv storage, [SQLite3][4] as log storage

## Dependencies

```sh
sudo apt install -y sqlite3 libsqlite3-dev librocksdb-dev libz-dev libbz2-dev libsnappy-dev liblz4-dev libzstd-dev
```

## Build and run

```sh
git clone https://github.com/csugulo/ShortLinkServer.git
cd ShortLinkServer

go run main.go
```

## Usage

```sh
# create a short url
curl --location --request POST 'localhost:8080/create' --data-raw '{"url":"https://www.bilibili.com/video/BV1Sx411T7QQ"}'

# show statistics
curl --location --request POST 'localhost:8080/statistics'
```

## Use Docker image
```
docker run -d -p 80:80 \
    --volume=$PWD/rocksdb:/opt/short_link_server/build/rocksdb \
    --volume=$PWD/sqlite.db:/opt/short_link_server/build/sqlite.db \
    csugulo/short_link_server -d YOUR_DOMAIN_NAME -p 80
```

## Try it online!
```
curl --location --request POST 'treemonkey.fun/create' --data-raw '{"url":"https://www.bilibili.com/video/BV1Sx411T7QQ"}'
```

[1]: https://go.dev/
[2]: https://github.com/gin-gonic/gin
[3]: https://github.com/facebook/rocksdb
[4]: https://www.sqlite.org/index.html