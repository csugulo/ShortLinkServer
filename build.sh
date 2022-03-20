#/bin/bash
rm -rf output
mkdir -p output

export CGO_CFLAGS="-I/usr/include"
export CGO_LDFLAGS="-L/usr/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd"
go build -o output/server main.go
cp -r pages output
cp bootstrap.sh output
cp config.yaml output