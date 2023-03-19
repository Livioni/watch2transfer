# watch2transfer

 transfer the inference results from edge to the cloud

## server

go run server/server.go

## client (edge device):

go run monitor/monitor.go

- input file_dir to be monitored
- input server's ip:port

## build for macOS

CGO_ENABLED=**0** GOOS=darwin GOARCH=amd64 go build -o mmonitor/monitor.go

## build for jetson

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o jetson monitor/monitor.go
