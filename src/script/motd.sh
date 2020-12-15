#!/bin/bash

ENV_FILE=`dirname "$(realpath "$0")"`/../../.env
if [ -f $ENV_FILE ]; then ENV_CONTENTS=`echo $(cat $ENV_FILE | sed 's/#.*//g' | xargs | grep serviceNames)`; fi
ENV_PARS=${ENV_CONTENTS##*=}
OLD_IFS=$IFS;IFS=,;SERVICES=($ENV_PARS);IFS=$OLD_IFS
upSeconds="$(/usr/bin/cut -d. -f1 /proc/uptime)"
secs=$((${upSeconds}%60))
mins=$((${upSeconds}/60%60))
hours=$((${upSeconds}/3600%24))
days=$((${upSeconds}/86400))
UPTIME=`printf "%d days, %02dh%02dm%02ds" "$days" "$hours" "$mins" "$secs"`
TotalHDD=`df -h | grep xvda1 | awk '{print $2}'`
UsingHDD=`df -h | grep xvda1 | awk '{print $3}'`
UsingHDDPer=`df -h | grep xvda1 | awk '{print $5}'`
PHP_FPM=`ps -A | awk '/php-fpm/{print $1}'`
NGINX=`ps -A | awk '/nginx/{print $1}'`
SOCKET=`ps aux | grep -v grep | grep httpd/Socket`
WEB_SOCKET=`ps aux | grep -v grep | grep httpd/WebSocket`
EMAIL_SENDER=`ps aux | grep -v grep | grep "/usr/bin/php /home/httpd/Apps/Cli/Run.php daemon run"`
EMAIL_LOGGER=`ps aux | grep -v grep | grep "/usr/bin/php /home/httpd/Apps/Cli/Run.php logger run"`

# get the load averages
read one five fifteen rest < /proc/loadavg

function in_array() {
    local needle array value
    needle="${1}"; shift; array=("${@}")
    for value in ${array[@]}; do [[ "${value}" =~ "${needle}" ]] && echo "true" && return; done
    echo "false"
}

echo "$(tput setaf 1)
`date +"%A, %e %B %Y, %r"`
`uname -srmo`
`/usr/bin/php -v | grep cli | awk '{print $1,$2}'`
Phalcon `/usr/bin/php -r "echo Phalcon\Version::get();"`$(tput setaf 2)
Location...........: `if [ -z $APP_ENV ]; then echo "production"; else echo $APP_ENV; fi`
Services...........: `for SERVICE in ${SERVICES[@]}; do printf "%s " ${SERVICE} | tr [a-z] [A-Z]; done`$(tput setaf 3)
Uptime.............: ${UPTIME}
Memory.............: `cat /proc/meminfo | grep MemFree | awk '{print $2}'`kB (Free) / `cat /proc/meminfo | grep MemTotal | awk '{print $2}'`kB (Total)
Load Averages......: ${one}, ${five}, ${fifteen} (1, 5, 15 min)
Private IP.........: `hostname -I`
Public IP..........: `curl -s http://checkip.amazonaws.com`
Disk...............: ${UsingHDD} / ${TotalHDD} (${UsingHDDPer})
PHP-FPM............: `if [[ -z $PHP_FPM ]]; then echo "$(tput setaf 1)OFF$(tput setaf 3)"; else echo "ON"; fi`
NGINX..............: `if [[ -z $NGINX ]]; then echo "$(tput setaf 1)OFF$(tput setaf 3)"; else echo "ON"; fi`
SOCKET.............: `if [[ -z $SOCKET ]]; then echo "$(tput setaf 1)OFF$(tput setaf 3)"; else echo "ON"; fi`
WEBSOCKET..........: `if [[ -z $WEB_SOCKET ]]; then echo "$(tput setaf 1)OFF$(tput setaf 3)"; else echo "ON"; fi`
EMAIL SENDER.......: `if [[ -z $EMAIL_SENDER ]]; then echo "$(tput setaf 1)OFF$(tput setaf 3)"; else echo "ON"; fi`
EMAIL LOGGER.......: `if [[ -z $EMAIL_LOGGER ]]; then echo "$(tput setaf 1)OFF$(tput setaf 3)"; else echo "ON"; fi`
$(tput sgr0)"