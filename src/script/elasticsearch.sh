#!/bin/bash
. ~/.bashrc
DAEMON=/usr/local/elasticsearch/jdk/bin/java
ES=/usr/local/elasticsearch/bin/elasticsearch

RET_VAL=0

check_process() {
	echo "elasticsearch shutdown check progress.."
	for i in {1..60}; do
		PID=($(pgrep -f ${DAEMON} | grep -v ^$$\$))
		if [[ -z ${PID} ]]; then
			echo "elasticsearch shutdown ok"
			RET_VAL=1
			return 1
		fi
		sleep 1
	done
	echo "ERROR: elasticsearch shutdown failed"
	RET_VAL=$?
}

daemon_start() {
	if [[ ! -f ${ES} ]]; then
		echo "elasticsearch not installed"
		RET_VAL=1
		return 1
	fi
	PID=($(pgrep -f ${DAEMON} | grep -v ^$$\$))
	if [[ ! -z ${PID} ]]; then
		echo "elasticsearch is already running."
		RET_VAL=1
		return 1
	fi
	${ES} -d > /dev/null 2>&1
	echo "Starting elasticsearch."
	RET_VAL=$?
}

daemon_stop() {
	PID=($(pgrep -f ${DAEMON} | grep -v ^$$\$))
	if [[ -z ${PID} ]]; then
		echo "elasticsearch is not running."
		RET_VAL=1
		return 1
	fi
	sudo kill ${PID}
	echo "Shutdown elasticsearch."
	RET_VAL=$?
}

case "$1" in
	start)
		daemon_start
		;;

	stop)
		daemon_stop
		;;

	restart|reload)
		daemon_stop
		check_process
		daemon_start
		RET_VAL=$?
		;;

	status)
		PID=($(pgrep -f $DAEMON | grep -v ^$$\$))
		if [[ -z ${PID} ]]; then
			echo "elasticsearch stopped"
		else
			echo "elasticsearch running. PID:" ${PID}
		fi
		;;
	*)
		echo "Usage: $0 {start|stop|restart|status}"
		exit 1
esac

exit ${RET_VAL}