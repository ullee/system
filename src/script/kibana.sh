#!/bin/bash
. ~/.bashrc
DAEMON=/usr/local/kibana/bin/../node/bin/node
KIBANA=/usr/local/kibana/bin/kibana

RET_VAL=0

check_process() {
	echo "kibana shutdown check progress.."
	for i in {1..60}; do
		PID=($(pgrep -f $DAEMON | grep -v ^$$\$))
		if [[ -z ${PID} ]]; then
			echo "kibana shutdown ok"
			RET_VAL=1
			return 1
		fi
		sleep 1
	done
	echo "ERROR: kibana shutdown failed"
	RET_VAL=$?
}

daemon_start() {
	if [[ ! -f ${KIBANA} ]]; then
		echo "kibana not installed"
		RET_VAL=1
		return 1
	fi
	PID=($(pgrep -f ${DAEMON} | grep -v ^$$\$))
	if [[ ! -z ${PID} ]]; then
		echo "kibana is already running."
		RET_VAL=1
		return 1
	fi
	${KIBANA} > /dev/null 2>&1 &
	echo "Starting kibana."
	RET_VAL=$?
}

daemon_stop() {
	PID=($(pgrep -f ${DAEMON} | grep -v ^$$\$))
	if [[ -z ${PID} ]]; then
		echo "kibana is not running."
		RET_VAL=1
		return 1
	fi
	sudo kill ${PID}
	echo "Shutdown kibana."
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
			echo "kibana stopped"
		else
			echo "kibana running. PID:" ${PID}
		fi
		;;
	*)
		echo "Usage: $0 {start|stop|restart|status}"
		exit 1
esac

exit ${RET_VAL}