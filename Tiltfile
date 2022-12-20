load('ext://helm_resource', 'helm_resource', 'helm_repo')

helm_repo('localstack', 'https://helm.localstack.cloud')
helm_resource('aws-local', 'localstack/localstack', port_forwards=[4566, 4566])