golang v1.14.5 download
```
wget https://golang.org/dl/go1.14.5.linux-amd64.tar.gz
```

golang v1.14.5 install
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
# GOPATH는 프로젝트 환경에 맞게 설정 한다. 의존성 패키지들이 해당 경로로 설치 된다.
export GOPATH=/home/system
source ~/.bashrc
```

build
```
cd $GOPATH/src/watchdog && go install
cd $GOPATH/src/watchdog-server && go install
```

execute server
```
$GOPATH/bin/watchdog-server
```

execute client
```
$GOPATH/bin/watchdog
```
