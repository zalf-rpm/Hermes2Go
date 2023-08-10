git describe --always --tags --long > version.txt
set /p VERSION=<version.txt
go build -v -ldflags "-X main.version=%VERSION%"