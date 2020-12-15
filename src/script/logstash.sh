#!/bin/bash
. ~/.bashrc
LOGSTASH_ENV=/home/httpd-log/Apps/Config/.bash_elk
if [[ -f ${LOGSTASH_ENV} ]]; then
    . ${LOGSTASH_ENV}
fi

DAEMON=logstash-core
LOGSTASH=/usr/local/logstash/bin/logstash
CONF=/home/httpd-log/Apps/Config/logstash.conf

RET_VAL=0

check_process() {
	echo "logstash shutdown check progress.."
	for i in {1..60}; do
		PID=($(pgrep -f ${DAEMON} | grep -v ^$$\$))
		if [[ -z ${PID} ]]; then
			echo "logstash shutdown ok"
			RET_VAL=1
			return 1
		fi
		sleep 1
	done
	echo "ERROR: logstash shutdown failed"
	RET_VAL=$?
}

daemon_start() {
	if [[ ! -f ${LOGSTASH} ]]; then
		echo "logstash not installed"
		RET_VAL=1
		return 1
	fi
	PID=($(pgrep -f ${DAEMON} | grep -v ^$$\$))
	if [[ ! -z ${PID} ]]; then
		echo "logstash is already running."
		RET_VAL=1
		return 1
	fi
	${LOGSTASH} -f ${CONF} > /dev/null 2>&1 &
	echo "Starting logstash."
	RET_VAL=$?
}

daemon_stop() {
	PID=($(pgrep -f ${DAEMON} | grep -v ^$$\$))
	if [[ -z ${PID} ]]; then
		echo "logstash is not running."
		RET_VAL=1
		return 1
	fi
	sudo kill ${PID}
	echo "Shutdown logstash."
	RET_VAL=$?
}

case "$1" in
	start)
		daemon_start
		;;

	stop)
		daemon_stop
		;;

	restart)
		daemon_stop
		check_process
		daemon_start
		RET_VAL=$?
		;;

	status)
		PID=($(pgrep -f ${DAEMON} | grep -v ^$$\$))
		if [[ -z ${PID} ]]; then
			echo "logstash stopped"
		else
			echo "logstash running. PID:" ${PID}
		fi
		;;
	*)
		echo "Usage: $0 {start|stop|restart|status}"
		exit 1
esac

exit ${RET_VAL}