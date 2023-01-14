.PHONY: run
run: 
	tilt down && tilt up --port 5050 

install:
	go get ./... && go mod tidy
build-workflow-engine:
	GOOS=linux go build  -o ./bin/workflow-server ./cmd/workflow-server
run-workflow-engine: build-workflow-engine
	ENDPOINT="http://localhost:4566" S3_FORCE_PATH_STYLE=true AWS_REGION="eu-west-2" STATE_MACHINE_ARN="arn:aws:states:eu-west-2:000000000000:stateMachine:notifications-workflow" ./bin/workflow-server
