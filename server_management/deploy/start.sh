#!/bin/bash

# Import common configurations
source ../common/config.sh

# Import deploy configurations
source config.sh

# Start the application
echo "Starting $MY_APP_NAME..."
cd $DEPLOY_DIR
nohup java -jar *.jar >$LOG_DIR/startup.log 2>&1 &
echo "Started $MY_APP_NAME"
