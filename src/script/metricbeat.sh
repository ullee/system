#!/bin/bash
. ~/.bashrc
DAEMON=/usr/local/metricbeat/metricbeat
METRICBEAT=/usr/local/metricbeat/metricbeat
METRICBEAT_CONF_DIR=/usr/local/metricbeat

RET_VAL=0

check_process() {
	echo "metricbeat shutdown check progress.."
	for i in {1..60}; do
		PID=($(pgrep -f $DAEMON | grep -v ^$$\$))
		if [[ -z ${PID} ]]; then
			echo "metricbeat shutdown ok"
			RET_VAL=1
			return 1
		fi
		sleep 1
	done
	echo "ERROR: metricbeat shutdown failed"
	RET_VAL=$?
}

daemon_start() {
	if [[ ! -f ${METRICBEAT} ]]; then
		echo "metricbeat not installed"
		RET_VAL=1
		return 1
	fi
	PID=($(pgrep -f ${DAEMON} | grep -v ^$$\$))
	if [[ ! -z ${PID} ]]; then
		echo "metricbeat is already running."
		RET_VAL=1
		return 1
	fi
	${METRICBEAT} -path.config ${METRICBEAT_CONF_DIR} > /dev/null 2>&1 &
	echo "Starting metricbeat."
	RET_VAL=$?
}

daemon_stop() {
	PID=($(pgrep -f ${DAEMON} | grep -v ^$$\$))
	if [[ -z ${PID} ]]; then
		echo "metricbeat is not running."
		RET_VAL=1
		return 1
	fi
	sudo kill ${PID}
	echo "Shutdown metricbeat."
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
			echo "metricbeat stopped"
		else
			echo "metricbeat running. PID:" ${PID}
		fi
		;;
	*)
		echo "Usage: $0 {start|stop|restart|status}"
		exit 1
esac

exit ${RET_VAL}