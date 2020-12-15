package main

import (
	. "constants"
	"custom-pkg/aws/cw"
	. "custom-pkg/logger"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Process struct {
	names  map[string]map[string]string
	consts Consts
}

func (p *Process) setProcess(service string) {

	var check string
	var start string

	switch service {
	case ES:
		check = "/usr/local/elasticsearch/jdk/bin/java"
		start = build.Default.GOPATH + "/bin/elasticsearch start"
	case LOGSTASH:
		check = "logstash-core"
		start = build.Default.GOPATH + "/bin/logstash start"
	case KIBANA:
		check = "/usr/local/kibana/bin/../node/bin/node"
		start = build.Default.GOPATH + "/bin/kibana start"
	case FILEBEAT:
		check = "/usr/local/filebeat/filebeat"
		start = build.Default.GOPATH + "/bin/filebeat start"
	case METRICBEAT:
		check = "/usr/local/metricbeat"
		start = build.Default.GOPATH + "/bin/metricbeat start"
	case PHP_FPM:
		check = "/var/run/php-fpm.pid"
		start = "sudo /usr/local/php/sbin/php-fpm -D"
	case NGINX:
		check = "/var/run/nginx.pid"
		start = "sudo /usr/local/nginx/sbin/nginx"
	}
	if _, ok := p.names[service]; !ok {
		p.names[service] = make(map[string]string)
		p.names[service]["check"] = check
		p.names[service]["start"] = start
	}
}

func (p *Process) runCommand(command string) error {
	commandSplit := strings.Fields(command)
	if len(commandSplit) < 2 {
		if _, err := exec.Command(command).Output(); err != nil {
			Log.Error(err)
			return err
		}
	} else {
		_, err := exec.Command(commandSplit[0], commandSplit[1:]...).Output()
		if err != nil {
			Log.Error(err)
			return err
		}
	}
	return nil
}

func (p *Process) check() {

	for processName, option := range p.names {

		pidStr := ""
		pid := 0

		switch processName {
		case ES, LOGSTASH, KIBANA, FILEBEAT, METRICBEAT:
			bytes, _ := exec.Command("pgrep", "-f", option["check"]).Output()
			pidStr = strings.Trim(string(bytes), "\n")

		default:
			if _, err := os.Stat(option["check"]); !os.IsNotExist(err) {
				bytes, err := ioutil.ReadFile(option["check"])
				if err != nil {
					Log.Error(processName, err)
					continue
				}
				pidStr = strings.Trim(string(bytes), "\n")
			}
		}

		if pidStr != "" {
			pidInt, err := strconv.Atoi(pidStr)
			if err != nil {
				Log.Error(processName, err)
				continue
			}
			pid = pidInt

			if _, err := os.FindProcess(pid); err != nil {
				fmt.Println(processName, err)
				Log.Error(processName, err)
				if err := p.runCommand(option["start"]); err == nil {
					fmt.Println(processName, "start ok")
					Log.Info(processName, "start ok")
				}
				continue
			}

			fmt.Println(processName, "is running")
			//Log.Debug(processName, "is running")

		} else {
			fmt.Println(processName, "is stopped")
			Log.Info(processName, "is stopped")
			if err := p.runCommand(option["start"]); err == nil {
				fmt.Println(processName, "start ok")
				Log.Info(processName, "start ok")
			}
		}

	}
}

func (p *Process) setCloudWatch() {

	metadata, err := cw.GetInstanceMetadata()
	if err != nil {
		Log.Panic(err)
	}

	memUtil, memUsed, memAvail, swapUtil, swapUsed, err := memoryUsage()

	var metricData []*cloudwatch.MetricDatum
	dims := cw.GetDimensions(metadata)

	metricData, err = cw.AddMetric("MemoryUtilization", "Percent", memUtil, dims, metricData)
	if err != nil {
		Log.Panic(err)
	}
	metricData, err = cw.AddMetric("MemoryUsed", "Bytes", memUsed, dims, metricData)
	if err != nil {
		Log.Panic(err)
	}
	metricData, err = cw.AddMetric("MemoryAvail", "Bytes", memAvail, dims, metricData)
	if err != nil {
		Log.Panic(err)
	}
	metricData, err = cw.AddMetric("SwapUsed", "Bytes", swapUsed, dims, metricData)
	if err != nil {
		Log.Panic(err)
	}
	metricData, err = cw.AddMetric("SwapUtil", "Percent", swapUtil, dims, metricData)
	if err != nil {
		Log.Panic(err)
	}

	path := "/"
	diskspaceUtil, diskspaceUsed, diskspaceAvail, diskinodesUtil, err := DiskSpace(path)
	if err != nil {
		Log.Panic(err)
	}
	metadata["fileSystem"] = path
	dims = cw.GetDimensions(metadata)
	metricData, err = cw.AddMetric("DiskUtilization", "Percent", diskspaceUtil, dims, metricData)
	if err != nil {
		Log.Panic(err)
	}
	metricData, err = cw.AddMetric("DiskUsed", "Bytes", float64(diskspaceUsed), dims, metricData)
	if err != nil {
		Log.Panic(err)
	}
	metricData, err = cw.AddMetric("DiskAvail", "Bytes", float64(diskspaceAvail), dims, metricData)
	if err != nil {
		Log.Panic(err)
	}
	metricData, err = cw.AddMetric("DiskInodesUtilization", "Percent", diskinodesUtil, dims, metricData)
	if err != nil {
		Log.Panic(err)
	}
	err = cw.PutMetric(metricData, "Linux/System")
	if err != nil {
		Log.Panic("Can't put CloudWatch Metric", err)
	}
}

func main() {

	var p = Process{}
	doc, _ := json.MarshalIndent(p.consts.ServiceNames(), "", "    ")
	service := flag.String("service", "", "ex) ./watchdog -service=nginx,php-fpm,es\nservices : "+string(doc))
	flag.Parse()

	if flag.NFlag() > 0 {
		p.names = map[string]map[string]string{}
		svc := strings.Split(*service, ",")
		for _, v := range svc {
			if _, ok := p.consts.ServiceNames()[v]; !ok {
				fmt.Println(v, "is invalid service name")
				Log.Fatal(v, "is invalid service name")
			}
			p.setProcess(v)
		}
		p.check()
	}
	p.setCloudWatch()
}
