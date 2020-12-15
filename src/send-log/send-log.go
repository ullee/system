package main

import (
	. "constants"
	"custom-pkg/aws/s3"
	. "custom-pkg/logger"
	"flag"
	"os"
	"os/exec"
	"sort"
	"syscall"
	"time"
)

type Context struct {
	baseDate   time.Time
	hostName   string
	s3         s3.Context
	consts     Consts
	files      []string
	service    string
	dir        string
	localPath  string
	uploadPath string
}

func (c *Context) send(isRootPermission ...bool) error {
	var err error
	c.s3.SetFileDir(c.localPath)
	c.s3.SetUploadDir(c.uploadPath)
	if err := c.s3.Upload(); err != nil {
		return err
	}
	if len(isRootPermission) > 0 && isRootPermission[0] == true {
		_, err := exec.Command("sudo", "unlink", c.localPath).Output()
		if err != nil {
			return err
		}
	} else {
		if err := syscall.Unlink(c.localPath); err != nil {
			return err
		}
	}
	Log.Debug("upload success", c.localPath)
	return err
}

func (c *Context) run() {
	for service, dir := range c.consts.LogPaths() {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			file, err := os.Open(dir)
			if err != nil {
				Log.Error(err)
				continue
			}
			files, _ := file.Readdirnames(0)
			_ = file.Close()
			if len(files) == 0 {
				Log.Debug(service, "log empty")
				continue
			}

			c.dir = dir
			c.service = service
			c.files = files
			sort.Sort(sort.StringSlice(c.files))

			switch c.service {
			case NGINX:
				c.nginx()
				break
			case ES:
				c.es()
				break
			case ES_GC:
				c.esGc()
				break
			case LOGSTASH:
				c.logstash()
				break
			case KIBANA:
				c.kibana()
				break
			case SYSTEM:
				c.system()
				break
			case NPRO, MTSCO:
				c.npro()
				break
			default:
				c.storage()
				break
			}
		}
	}
}

func main() {

	days := flag.Int("days", 7, "days")
	flag.Parse()

	Log.Info("start send-log")
	if *days > 0 {
		*days *= -1
	}

	c := new(Context)
	now := time.Now()
	// 설정된 일자 이하 (ex. 7일을 설정하면 7일치만 남기고 나머지 로그 데이터를 정리함)
	c.baseDate = now.AddDate(0, 0, *days)

	c.hostName, _ = os.Hostname()
	Log.Info("baseDate:", c.baseDate)
	c.run()

	Log.Info("end send-log")
}
