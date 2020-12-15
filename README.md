download
```
wget https://golang.org/dl/go1.14.5.linux-amd64.tar.gz
```

install
```
tar xzvBpf go1.14.5.linux-amd64.tar.gz && \
sudo mv go /usr/local/src/go1.14.5 && \
sudo ln -s /usr/local/src/go1.14.5 /usr/local/go && \
sudo ln -s /usr/local/go/bin/* /usr/bin/ && \
unlink go1.14.5.linux-amd64.tar.gz
```

env
```
vi ~/.bashrc
export GOPATH=/home/system
source ~/.bashrc
```

build
```
cd /home/system/src/watchdog && go install
```

execute
```
/home/system/bin/watchdog
```
