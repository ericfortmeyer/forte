#!/bin/sh

trap 'kill $PID' SIGTERM
php -t /srv/php-app -S 0.0.0.0:8000 &
PID=$!
wait $PID
trap - SIGTERM
wait $PID
