#!/bin/bash

# Import monitoring configurations
source config.sh

# Send alert message
echo "$1" | mail -s "[ALERT] $MY_APP_NAME" $ALERT_EMAIL
