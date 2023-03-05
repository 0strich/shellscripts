#!/bin/bash

# Import common configurations
source ../common/config.sh

# Import deploy configurations
source config.sh

# Deploy the application
echo "Deploying $MY_APP_NAME..."
cd $MY_APP_HOME
git pull
./gradlew build
cp build/libs/*.jar $DEPLOY_DIR
echo "Deployed $MY_APP_NAME to $DEPLOY_DIR"
