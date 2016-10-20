#! bin/sh

rm -rf dist &&
mkdir dist &&

mkdir dist/linux/ &&
cp -rf public dist/linux &&
cp app-config.toml dist/linux/app-config.toml &&
cp sql/init.sql dist/linux/init.sql &&
env GOOS=linux GOARCH=386 go build -o dist/linux/tuprt_32 main.go &&
env GOOS=linux GOARCH=amd64 go build -o dist/linux/tuprt_64 main.go &&

mkdir dist/mac/ &&
cp -rf public dist/mac &&
cp app-config.toml dist/mac/app-config.toml &&
cp sql/init.sql dist/mac/init.sql &&
env GOOS=darwin GOARCH=386 go build -o dist/mac/tuprt_32 main.go &&
env GOOS=darwin GOARCH=amd64 go build -o dist/mac/tuprt_64 main.go &&

mkdir dist/windows/ &&
cp -rf public dist/windows &&
cp app-config.toml dist/windows/app-config.toml &&
cp sql/init.sql dist/windows/init.sql &&
env GOOS=windows GOARCH=386 go build -o dist/windows/tuprt_32.exe main.go &&
env GOOS=windows GOARCH=amd64 go build -o dist/windows/tuprt_64.exe main.go