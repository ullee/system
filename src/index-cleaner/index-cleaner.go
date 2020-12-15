package main

import (
	"bytes"
	. "constants"
	. "custom-pkg/logger"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type DefaultResult struct {
	Acknowledged bool `json:"acknowledged"`
}

type SnapshotResult struct {
	Snapshot struct {
		Snapshot           string   `json:"snapshot"`
		Uuid               string   `json:"uuid"`
		VersionId          int      `json:"version_id"`
		Version            string   `json:"version"`
		Indices            []string `json:"indices"`
		IndicesGlobalState bool     `json:"indices_global_state"`
		State              string   `json:"state"`
		StartTime          string   `json:"start_time"`
		StartTimeInMillis  uint64   `json:"start_time_in_millis"`
		EndTime            string   `json:"end_time"`
		EndTimeInMillis    uint64   `json:"end_time_in_millis"`
		DurationInMillis   uint64   `json:"duration_in_millis"`
		Failures           []string `json:"failures"`
		Shards             struct {
			Total      int `json:"total"`
			Failed     int `json:"failed"`
			Successful int `json:"successful"`
		} `json:"shards"`
	} `json:"snapshot"`
}

type Repository struct {
	Type     string   `json:"type"`
	Settings Settings `json:"settings"`
}

type Settings struct {
	Bucket    string `json:"bucket"`
	BasePath  string `json:"base_path"`
	Region    string `json:"region"`
	Compress  bool   `json:"compress"`
	CannedAcl string `json:"canned_acl"`
	Readonly  bool   `json:"readonly"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type Snapshots struct {
	Indices            string `json:"indices"`
	IgnoreUnavailable  bool   `json:"ignore_unavailable"`
	IncludeGlobalState bool   `json:"include_global_state"`
}

func request(method string, url string, body io.Reader) string {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		Log.Error(err)
		return ""
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		Log.Error(err)
		return ""
	}
	defer response.Body.Close()

	bytes, _ := ioutil.ReadAll(response.Body)
	return string(bytes)
}

func main() {

	host := flag.String("host", "127.0.0.1", "Elasticsearch host ip address")
	port := flag.String("port", "9200", "Elasticsearch port")
	days := flag.Int("days", 90, "Number of days to leave behind Elasticsearch index")

	flag.Parse()

	/*
		if flag.NFlag() == 0 {
			flag.Usage()
			return
		}
	*/

	Log.Info("start index-cleaner for Elasticsearch")
	Log.Info("leave index. days :", *days)

	if *days > 0 {
		*days *= -1
	}

	url := "http://" + *host + ":" + *port + "/_cluster/health?level=indices"
	response := request("GET", url, nil)

	if response == "" {
		Log.Panic("response is null. url", url)
	}

	data := make(map[string]interface{}) // response json data를 담기 위한 map 초기화

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		Log.Panic(err)
	}

	indices := make([]string, 0) // indices 슬라이스 생성
	for _, value := range data {
		if indexes, isObj := value.(map[string]interface{}); isObj {
			for index, _ := range indexes {
				indices = append(indices, index)
			}
		} else {
			// do nothing.
		}
	}

	now := time.Now() // 현재 날짜
	baseDate := now.AddDate(0, 0, *days)
	Log.Info("baseDate:", baseDate)

	var bucket string
	if os.Getenv("APP_ENV") == "production" {
		bucket = S3_BUCKET_PRODUCTION
	} else {
		bucket = S3_BUCKET_STAGING
	}

	var result DefaultResult
	var sresult SnapshotResult

	for _, indice := range indices {
		indiceSlice := strings.Split(indice, "-")
		// 시스템 인덱스 정리
		if len(indiceSlice) > 3 && indiceSlice[0] == ".monitoring" {
			utc, _ := time.Parse("2006.01.02", indiceSlice[3])
			loc, _ := time.LoadLocation("Asia/Seoul")
			kst := utc.In(loc)
			if kst.Before(baseDate) {
				url := "http://" + *host + ":" + *port + "/" + indice
				response = request("DELETE", url, nil)
				//response := `{"acknowledged":true}`

				if response == "" {
					Log.Error("response is null. url:", url)
					continue
				}

				if err := json.Unmarshal([]byte(response), &result); err != nil {
					Log.Error(err)
				}
				if result.Acknowledged {
					Log.Info("delete success:", url)
				}
			}
			// 인덱스 S3 백업
		} else if len(indiceSlice) > 1 && indiceSlice[0] == "gateway" {
			repoName := bucket + "-gateway"
			utc, _ := time.Parse("2006.01.02", indiceSlice[1])
			loc, _ := time.LoadLocation("Asia/Seoul")
			kst := utc.In(loc)
			if kst.Before(baseDate) {
				url = "http://" + *host + ":" + *port + "/_snapshot/" + repoName + "?pretty"
				ok := request("GET", url, nil)
				if ok != "" {
					repository := Repository{
						Type: "s3",
						Settings: struct {
							Bucket    string `json:"bucket"`
							BasePath  string `json:"base_path"`
							Region    string `json:"region"`
							Compress  bool   `json:"compress"`
							CannedAcl string `json:"canned_acl"`
							Readonly  bool   `json:"readonly"`
							AccessKey string `json:"access_key"`
							SecretKey string `json:"secret_key"`
						}{
							Bucket:    bucket,
							BasePath:  "es-indices-backup/gateway",
							Region:    "ap-northeast-2",
							Compress:  true,
							CannedAcl: "public-read",
							Readonly:  false,
							AccessKey: S3_ACCESS_KEY,
							SecretKey: S3_SECRET_KEY,
						},
					}
					pBytes, _ := json.Marshal(repository)
					buff := bytes.NewBuffer(pBytes)
					response = request("PUT", url, buff)

					if err := json.Unmarshal([]byte(response), &result); err != nil {
						Log.Panic(err)
						return
					}

					if result.Acknowledged {
						Log.Info("create snapshot repository complete.", url)
					}
				}

				url = "http://" + *host + ":" + *port + "/_snapshot/" + repoName + "/" + indice + "?wait_for_completion=true"
				snapshot := Snapshots{
					Indices:            indice,
					IgnoreUnavailable:  true,
					IncludeGlobalState: false,
				}
				pBytes, _ := json.Marshal(snapshot)
				buff := bytes.NewBuffer(pBytes)
				response = request("PUT", url, buff)
				if err := json.Unmarshal([]byte(response), &sresult); err != nil {
					Log.Error(err)
					continue
				}
				if sresult.Snapshot.Shards.Failed > 0 {
					Log.Error("Create Snapshot failed..", url)
					continue
				}

				url := "http://" + *host + ":" + *port + "/" + indice
				response = request("DELETE", url, nil)
				//response := `{"acknowledged":true}`

				if response == "" {
					Log.Error("response is null. url:", url)
					continue
				}

				if err := json.Unmarshal([]byte(response), &result); err != nil {
					Log.Error(err)
				}
				if result.Acknowledged {
					Log.Info("delete success:", url)
				}
			}
		}
	}

	Log.Info("end index-cleaner")
}
