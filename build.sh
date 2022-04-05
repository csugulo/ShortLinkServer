#/bin/bash
rm -rf build
mkdir -p build

export CGO_CFLAGS="-I/usr/include"
export CGO_LDFLAGS="-L/usr/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd"
go build -o build/server main.go