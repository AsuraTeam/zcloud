export GOPATH=//home/gitsrc/cloud/
export GOROOT=/usr/local/go
cd  /home/gitsrc/cloud/src/cloud/
go build -i -o  cloud-server 
ps aux |grep cloud-server |grep -v grep|awk '{print "kill -9 "$2}'|bash
sleep 9
./cloud-server
