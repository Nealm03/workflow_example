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
  deps=['./tilt/setup_workflow_locally.sh', './src/statemachine.json', './src/notifyUser.ts', './src/findUsersToNotify.ts'],
)