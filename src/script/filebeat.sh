#!/bin/bash
. ~/.bashrc
DAEMON=/usr/local/filebeat/filebeat
FILEBEAT=/usr/local/filebeat/filebeat
FILEBEAT_CONF_DIR=/usr/local/filebeat

RET_VAL=0

check_process() {
	echo "filebeat shutdown check progress.."
	for i in {1..60}; do
		PID=($(pgrep -f $DAEMON | grep -v ^$$\$))
		if [[ -z ${PID} ]]; then
			echo "filebeat shutdown ok"
			RET_VAL=1
			return 1
		fi
		sleep 1
	done
	echo "ERROR: filebeat shutdown failed"
	RET_VAL=$?
}

daemon_start() {
	if [[ ! -f ${FILEBEAT} ]]; then
		echo "filebeat not installed"
		RET_VAL=1
		return 1
	fi
	PID=($(pgrep -f ${DAEMON} | grep -v ^$$\$))
	if [[ ! -z ${PID} ]]; then
		echo "filebeat is already running."
		RET_VAL=1
		return 1
	fi
	${FILEBEAT} -path.config ${FILEBEAT_CONF_DIR} > /dev/null 2>&1 &
	echo "Starting filebeat."
	RET_VAL=$?
}

daemon_stop() {
	PID=($(pgrep -f ${DAEMON} | grep -v ^$$\$))
	if [[ -z ${PID} ]]; then
		echo "filebeat is not running."
		RET_VAL=1
		return 1
	fi
	sudo kill ${PID}
	echo "Shutdown filebeat."
	RET_VAL=$?
}

case "$1" in
	start)
		daemon_start
		;;

	stop)
		daemon_stop
		check_process
		;;

	restart|reload)
		daemon_stop
		check_process
		daemon_start
		RET_VAL=$?
		;;

	status)
		PID=($(pgrep -f ${DAEMON} | grep -v ^$$\$))
		if [[ -z ${PID} ]]; then
			echo "filebeat stopped"
		else
			echo "filebeat running. PID:" ${PID}
		fi
		;;
	*)
		echo "Usage: $0 {start|stop|restart|status}"
		exit 1
esac

exit ${RET_VAL}