#!/bin/bash
DEFAULT_REGION="eu-west-2"
WORKFLOW_NAME="notifications-workflow"
DUMMY_STATE_MACHINE_ARN="arn:aws:states:$DEFAULT_REGION:000000000000:stateMachine:$WORKFLOW_NAME"

BASE_PATH=$PWD
SOURCE_CODE_WORKSPACE_DIR="./src"
INFRA_DIR="./infra"

function build_lambdas(){
  cd $SOURCE_CODE_WORKSPACE_DIR
  npm ci && npm run build
  cd ./dist
  zip -r ./findUsersToNotify.zip .
  zip -r ./notifyUser.zip .
  cd $BASE_PATH
}

function deploy_infra(){
  cd $INFRA_DIR
  cdklocal bootstrap
  cdklocal deploy --require-approval never
}

function invoke_workflow(){
  EXECUTION_ARN=$(awslocal stepfunctions start-execution \
    --region $DEFAULT_REGION  \
    --state-machine-arn $DUMMY_STATE_MACHINE_ARN \
    --name test \
    --input "{\"olderThan\": 20}" | jq -r .executionArn)
  sleep 5
  awslocal stepfunctions get-execution-history --execution-arn  $EXECUTION_ARN --region $DEFAULT_REGION | jq 
}

function teardown(){
  rm -rf $WORKSPACE_DIR/dist
  cd $INFRA_DIR
  cdklocal destroy --force
  cd $BASE_PATH
}

teardown
deploy_infra
invoke_workflow
