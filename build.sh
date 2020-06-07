env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o builds/geek-life_darwin-amd64 ./app
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/geek-life_linux-amd64 ./app
env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o builds/geek-life_linux-arm64 ./app
env GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o builds/geek-life_windows-386 ./app
upx builds/geek-life_*