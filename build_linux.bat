go build -ldflags="-w -s"
set GOARCH=amd64
set GOOS=linux
go build -ldflags="-w -s"