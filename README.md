# 4lance
Build for RaspberryPi 2
GOOS=linux GOARCH=arm GOARM=7 go build -o mainARMv7 main.go
access to remote mongo
ssh -L 4321:localhost:27017 user@your.ip.address -f -N
mongo --port 4321