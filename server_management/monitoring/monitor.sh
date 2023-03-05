#!/bin/bash

# Import common configurations
source ../common/config.sh

# Import monitoring configurations
source config.sh

# Check server status
if ! curl -sSf "http://localhost:$MY_APP_PORT" >/dev/null; then
	echo "Server is down! Restarting..."
	../restart.sh
	if ! curl -sSf "http://localhost:$MY_APP_PORT" >/dev/null; then
		echo "Restart failed! Sending alert..."
		../alert.sh "Server is down and could not be restarted!"
	fi
fi
