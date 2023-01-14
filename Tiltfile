load('ext://helm_resource', 'helm_resource', 'helm_repo')

helm_repo(
  'localstack', 
  'https://helm.localstack.cloud',
)

helm_resource(
  'aws-local', 
  'localstack/localstack', 
  port_forwards=[4566],
)

local_resource(
  'setup-local', 
  cmd='./tilt/setup_workflow_locally.sh', 
  deps=['./activities'],
)

local_resource(
  'workflow-engine', 
  cmd='make build-workflow-engine', 
  serve_env={
    'ENDPOINT':'http://localhost:4566', 
    'AWS_REGION':'eu-west-2',
    'STATE_MACHINE_ARN':'arn:aws:states:eu-west-2:000000000000:stateMachine:notifications-workflow'
  },
  serve_cmd='./bin/workflow-server',
  deps=['./cmd', './internal']
)