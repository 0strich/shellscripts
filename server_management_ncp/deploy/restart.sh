#!/bin/bash

# Import common configurations
source ../common/config.sh

# Import deploy configurations
source config.sh

# Restart the application
echo "Restarting $MY_APP_NAME..."
../stop.sh
../start.sh
echo "Restarted $MY_APP_NAME"
