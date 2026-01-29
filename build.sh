# go build -o geek-life ./app

echo "# building:"
env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o builds/geek-life_darwin-amd64 ./app
# env GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o builds/geek-life_darwin-arm64 ./app
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/geek-life_linux-amd64 ./app
env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o builds/geek-life_linux-arm64 ./app
env GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o builds/geek-life_windows-386 ./app

echo "running upx:"
if command -v upx >/dev/null 2>&1; then
  upx builds/geek-life_*
else
  echo "upx not found!"
fi

echo "SHA256 sum of release binaries:"
if command -v shasum >/dev/null 2>&1; then
  shasum -a 256 -b builds/geek-life_*
else
  echo "shasum not found!"
fi
