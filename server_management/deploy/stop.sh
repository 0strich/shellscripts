#!/bin/bash

# Import common configurations
source ../common/config.sh

# Import deploy configurations
source config.sh

# Stop the application
echo "Stopping $MY_APP_NAME..."
PID=$(ps aux | grep "$MY_APP_NAME" | grep -v grep | awk '{print $2}')
kill $PID
echo "Stopped $MY_APP_NAME"
