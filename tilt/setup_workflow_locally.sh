#!/bin/bash
WORKFLOW_NAME="NotifyUsersStateMachine"
DUMMY_LAMBDA_ROLE_ARN="arn:aws:iam::012345678901:role/DummyRole"
DUMMY_STATE_MACHINE_ARN="arn:aws:states:us-east-1:000000000000:stateMachine:$WORKFLOW_NAME"

WORKSPACE_DIR="./src"


function navigate_to_workspace_dir(){
  cd $WORKSPACE_DIR
}

function build_lambdas(){
  npm ci && npm run build
  cd ./dist
  zip -r ./findUsersToNotify.zip .
  zip -r ./notifyUser.zip .
  cd ..
}

function publish_lambdas(){
  awslocal lambda create-function --function-name FindUsersToNotify \
    --runtime nodejs14.x \
    --handler findUsersToNotify.findUsersToNotify \
    --role $DUMMY_LAMBDA_ROLE_ARN \
    --zip-file fileb://dist/findUsersToNotify.zip
}


function create_stepfunctions(){
  awslocal stepfunctions create-state-machine \
    --definition  file://statemachine.json \
    --name $WORKFLOW_NAME \
    --role-arn $DUMMY_LAMBDA_ROLE_ARN

}

function invoke_workflow(){
  awslocal stepfunctions start-execution --name $WORKFLOW_NAME \
    --state-machine-arn $DUMMY_STATE_MACHINE_ARN \
    --name test \
    --input "{\"olderThan\": 20}"
  sleep 5 
  awslocal stepfunctions list-executions --state-machine-arn $DUMMY_STATE_MACHINE_ARN
}

function teardown(){
  rm -rf $WORKSPACE_DIR/dist
  awslocal lambda delete-function --function-name FindUsersToNotify

  awslocal stepfunctions delete-state-machine --state-machine-arn $DUMMY_STATE_MACHINE_ARN
}

teardown
navigate_to_workspace_dir
build_lambdas
publish_lambdas
create_stepfunctions
invoke_workflow