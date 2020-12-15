package main

import (
	. "custom-pkg/logger"
	"strings"
	"syscall"
	"time"
)

func (c *Context) nginx() {
	for _, file := range c.files {
		temp := strings.Split(file, ".")
		if len(temp) > 2 && temp[2] == "gz" {
			fileNameSlice := strings.Split(temp[1], "-")
			if len(fileNameSlice) > 1 {
				year := fileNameSlice[1][:4]
				month := fileNameSlice[1][4:6]
				day := fileNameSlice[1][6:8]
				utc, _ := time.Parse("2006.01.02", year+"."+month+"."+day)
				loc, _ := time.LoadLocation("Asia/Seoul")
				kst := utc.In(loc)
				if kst.Before(c.baseDate) {
					c.localPath = c.dir + "/" + file
					c.uploadPath = "/logs/" + c.service + "/" + c.hostName + "/" + year + "/" + month + "/" + day + "/" + file
					if err := c.send(true); err != nil {
						Log.Error(err)
						continue
					}
				}
			}
		}
	}
}

func (c *Context) es() {
	for _, file := range c.files {
		temp := strings.Split(file, ".")
		if len(temp) > 2 && temp[2] == "gz" {
			fileNameSlice := strings.Split(temp[0], "-")
			if len(fileNameSlice) > 4 {
				year := fileNameSlice[1]
				month := fileNameSlice[2]
				day := fileNameSlice[3]
				utc, _ := time.Parse("2006.01.02", year+"."+month+"."+day)
				loc, _ := time.LoadLocation("Asia/Seoul")
				kst := utc.In(loc)
				if kst.Before(c.baseDate) {
					c.localPath = c.dir + "/" + file
					c.uploadPath = "/logs/" + c.service + "/" + c.hostName + "/" + year + "/" + month + "/" + day + "/" + file
					if err := c.send(); err != nil {
						Log.Error(err)
						continue
					}
				}
			}
		}
	}
}

func (c *Context) esGc() {
	for _, file := range c.files {
		if err := syscall.Unlink(c.dir + "/" + file); err != nil {
			Log.Error(err)
			continue
		}
	}
	Log.Info("delete success es gc log")
}

func (c *Context) logstash() {
	for _, file := range c.files {
		temp := strings.Split(file, ".")
		if len(temp) > 2 && temp[2] == "gz" {
			fileNameSlice := strings.Split(temp[0], "-")
			if len(fileNameSlice) > 5 {
				year := fileNameSlice[2]
				month := fileNameSlice[3]
				day := fileNameSlice[4]
				utc, _ := time.Parse("2006.01.02", year+"."+month+"."+day)
				loc, _ := time.LoadLocation("Asia/Seoul")
				kst := utc.In(loc)
				if kst.Before(c.baseDate) {
					c.localPath = c.dir + "/" + file
					c.uploadPath = "/logs/" + c.service + "/" + c.hostName + "/" + year + "/" + month + "/" + day + "/" + file
					if err := c.send(); err != nil {
						Log.Error(err)
						continue
					}
				}
			}
		}
	}
}

func (c *Context) kibana() {
	for _, file := range c.files {
		temp := strings.Split(file, ".")
		if len(temp) > 2 && temp[2] == "gz" {
			fileNameSlice := strings.Split(temp[0], "-")
			if len(fileNameSlice) > 1 {
				year := fileNameSlice[1][:4]
				month := fileNameSlice[1][4:6]
				day := fileNameSlice[1][6:8]
				utc, _ := time.Parse("2006.01.02", year+"."+month+"."+day)
				loc, _ := time.LoadLocation("Asia/Seoul")
				kst := utc.In(loc)
				if kst.Before(c.baseDate) {
					c.localPath = c.dir + "/" + file
					c.uploadPath = "/logs/" + c.service + "/" + c.hostName + "/" + year + "/" + month + "/" + day + "/" + file
					if err := c.send(); err != nil {
						Log.Error(err)
						continue
					}
				}
			}
		}
	}
}

func (c *Context) storage() {
	for _, file := range c.files {
		temp := strings.Split(file, ".")
		if len(temp) > 1 && temp[1] == "log" {
			fileNameSlice := strings.Split(temp[0], "-")
			if len(fileNameSlice) > 3 {
				year := fileNameSlice[1]
				month := fileNameSlice[2]
				day := fileNameSlice[3]
				utc, _ := time.Parse("2006.01.02", year+"."+month+"."+day)
				loc, _ := time.LoadLocation("Asia/Seoul")
				kst := utc.In(loc)
				if kst.Before(c.baseDate) {
					c.localPath = c.dir + "/" + file
					c.uploadPath = "/logs/storage/" + c.hostName + "/" + year + "/" + month + "/" + day + "/" + c.service + "/" + file
					if err := c.send(); err != nil {
						Log.Error(err)
						continue
					}
				}
			}
			// NOTIFY,LOG 데몬 로그
		} else if len(temp) > 2 && temp[2] == "gz" {
			fileNameSlice := strings.Split(temp[1], "-")
			if len(fileNameSlice) > 1 {
				year := fileNameSlice[1][:4]
				month := fileNameSlice[1][4:6]
				day := fileNameSlice[1][6:8]
				utc, _ := time.Parse("2006.01.02", year+"."+month+"."+day)
				loc, _ := time.LoadLocation("Asia/Seoul")
				kst := utc.In(loc)
				if kst.Before(c.baseDate) {
					c.localPath = c.dir + "/" + file
					c.uploadPath = "/logs/storage/" + c.hostName + "/" + year + "/" + month + "/" + day + "/" + c.service + "/" + file
					if err := c.send(true); err != nil {
						Log.Error(err)
						continue
					}
				}
			}
		}
	}
}

func (c *Context) system() {
	for _, file := range c.files {
		temp := strings.Split(file, ".")
		if len(temp) > 1 && temp[1] == "log" {
			fileNameSlice := strings.Split(temp[0], "-")
			if len(fileNameSlice) > 3 {
				year := fileNameSlice[1]
				month := fileNameSlice[2]
				day := fileNameSlice[3]
				utc, _ := time.Parse("2006.01.02", year+"."+month+"."+day)
				loc, _ := time.LoadLocation("Asia/Seoul")
				kst := utc.In(loc)
				if kst.Before(c.baseDate) {
					c.localPath = c.dir + "/" + file
					c.uploadPath = "/logs/" + c.service + "/" + c.hostName + "/" + year + "/" + month + "/" + day + "/" + file
					if err := c.send(); err != nil {
						Log.Error(err)
						continue
					}
				}
			}
		}
	}
}

func (c *Context) npro() {
	for _, file := range c.files {
		temp := strings.Split(file, ".")
		if len(temp) > 2 && temp[1] == "log" {
			year := temp[2][:4]
			month := temp[2][4:6]
			day := temp[2][6:8]
			utc, _ := time.Parse("2006.01.02", year+"."+month+"."+day)
			loc, _ := time.LoadLocation("Asia/Seoul")
			kst := utc.In(loc)
			if kst.Before(c.baseDate) {
				c.localPath = c.dir + "/" + file
				c.uploadPath = "/logs/" + c.service + "/" + c.hostName + "/" + year + "/" + month + "/" + day + "/" + file
				if err := c.send(true); err != nil {
					Log.Error(err)
					continue
				}
			}
		}
	}
}
