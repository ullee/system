package constants

const (
	NGINX      = "nginx"
	PHP_FPM    = "php-fpm"
	ES         = "es"
	ES_GC      = "es-gc"
	LOGSTASH   = "logstash"
	KIBANA     = "kibana"
	FILEBEAT   = "filebeat"
	METRICBEAT = "metricbeat"
	DEFAULT    = "default"
	GW         = "gw"
	BANNER     = "banner"
	BOARD      = "board"
	CURATOR    = "curator"
	DEAL       = "deal"
	DELIVERY   = "delivery"
	EVENT      = "event"
	EXTERNAL   = "external"
	SOCKET     = "socket"
	STORAGE    = "storage"
	TICKET     = "ticket"
	NOTIFY     = "notify"
	LOG        = "log"
	DEVOPS     = "devops"
	SYSTEM     = "system"
	NPRO       = "npro"
	MTSCO      = "mtsco"
)

const (
	TAG_FRONT  = "front"
	TAG_CORE   = "core"
	TAG_DEVOPS = "devops"
	TAG_ES     = "es"
	TAG_SOCKET = "socket"
	TAG_CACHE  = "cache"
)

const (
	DEFAULT_STAGING_CRITICAL_VALUE    = 99.0
	DEFAULT_PRODUCTION_CRITICAL_VALUE = 96.0
)

type Consts struct {
	LogPath       map[string]string
	ServiceName   map[string]string
	CriticalValue map[int]map[string]float64
}

func (c *Consts) ServiceNames() map[string]string {
	c.ServiceName = map[string]string{}
	c.ServiceName[NGINX] = "nginx"
	c.ServiceName[PHP_FPM] = "php-fpm"
	c.ServiceName[ES] = "es"
	c.ServiceName[LOGSTASH] = "logstash"
	c.ServiceName[KIBANA] = "kibana"
	c.ServiceName[FILEBEAT] = "filebeat"
	c.ServiceName[METRICBEAT] = "metricbeat"
	c.ServiceName[SYSTEM] = "system"
	c.ServiceName[NPRO] = "npro"
	c.ServiceName[MTSCO] = "mtsco"
	return c.ServiceName
}

func (c *Consts) LogPaths() map[string]string {
	c.LogPath = map[string]string{}
	c.LogPath[NGINX] = "/var/log/nginx"
	c.LogPath[ES] = "/var/log/elasticsearch"
	c.LogPath[ES_GC] = "/usr/local/elasticsearch/logs"
	c.LogPath[LOGSTASH] = "/var/log/logstash"
	c.LogPath[KIBANA] = "/var/log/kibana"
	c.LogPath[DEFAULT] = "/home/httpd/Storage/Logs"
	c.LogPath[GW] = "/home/httpd-api/Storage/Logs"
	c.LogPath[BANNER] = "/home/httpd-banner/Storage/Logs"
	c.LogPath[BOARD] = "/home/httpd-board/Storage/Logs"
	c.LogPath[CURATOR] = "/home/httpd-curator/Storage/Logs"
	c.LogPath[DEAL] = "/home/httpd-deal/Storage/Logs"
	c.LogPath[DELIVERY] = "/home/httpd-delivery/Storage/Logs"
	c.LogPath[EVENT] = "/home/httpd-event/Storage/Logs"
	c.LogPath[EXTERNAL] = "/home/httpd-external/Storage/Logs"
	c.LogPath[STORAGE] = "/home/httpd-storage/Storage/Logs"
	c.LogPath[TICKET] = "/home/httpd-ticket/Storage/Logs"
	c.LogPath[LOG] = "/home/httpd-log/Storage/Logs"
	c.LogPath[DEVOPS] = "/home/httpd-devops/Storage/Logs"
	c.LogPath[SOCKET] = "/home/httpd/logs"
	c.LogPath[SYSTEM] = "/home/system/logs"
	c.LogPath[NPRO] = "/var/log/sms/npro2"
	c.LogPath[MTSCO] = "/var/log/sms/mtsco"
	return c.LogPath
}

const (
	METRIC_TYPE_DISK   = 1
	METRIC_TYPE_CPU    = 2
	METRIC_TYPE_MEMORY = 3
)

func (c *Consts) CriticalValues() map[int]map[string]float64 {
	c.CriticalValue = map[int]map[string]float64{}
	// Disk
	c.CriticalValue[METRIC_TYPE_DISK] = map[string]float64{
		TAG_DEVOPS: 96.0,
		TAG_ES:     96.0,
		TAG_FRONT:  95.0,
		TAG_SOCKET: 98.0,
		TAG_CACHE:  96.0,
	}
	// CPU
	c.CriticalValue[METRIC_TYPE_CPU] = map[string]float64{
		TAG_DEVOPS: 90.0,
		TAG_ES:     65.0,
		TAG_FRONT:  95.0,
		TAG_SOCKET: 20.0,
		TAG_CACHE:  90.0,
	}
	// Memory
	c.CriticalValue[METRIC_TYPE_MEMORY] = map[string]float64{
		TAG_DEVOPS: 96.0,
		TAG_ES:     80.0,
		TAG_FRONT:  95.0,
		TAG_SOCKET: 30.0,
		TAG_CACHE:  99.0,
	}
	return c.CriticalValue
}

const (
	SLACK_WEBHOOKURL = "Slack webhook URL"
)

const (
	EC2_REGION     = "ap-northeast-2"
	EC2_ACCESS_KEY = "EC2_ACCESS_KEY"
	EC2_SECRET_KEY = "EC2_SECRET_KEY"

	SSM_REGION     = "ap-northeast-2"
	SSM_ACCESS_KEY = "SSM_ACCESS_KEY"
	SSM_SECRET_KEY = "SSM_SECRET_KEY"

	S3_REGION            = "ap-northeast-2"
	S3_BUCKET_STAGING    = "staging"
	S3_BUCKET_PRODUCTION = "production"
	S3_ACCESS_KEY        = "S3_ACCESS_KEY"
	S3_SECRET_KEY        = "S3_SECRET_KEY"

	CW_REGION     = "ap-northeast-2"
	CW_ACCESS_KEY = "CW_ACCESS_KEY"
	CW_SECRET_KEY = "CW_SECRET_KEY"
)
