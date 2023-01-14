package internal

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

type StartWorkflowCommand struct {
	WorkflowId   *string // optional,
	WorkflowName string  // required
	Input        JsonInput
}

type WorkflowApi interface {
	StartWorkflow(ctx context.Context, cmd StartWorkflowCommand) (workflowId string, err error)
	GetWorkflowStatus(ctx context.Context, workflowId string) (*WorkflowExecution, error)
}

type JsonInput *string

func (api *stepFunctionWorkflowApi) StartWorkflow(ctx context.Context, cmd StartWorkflowCommand) (string, error) {
	output, err := api.StartExecution(ctx, &sfn.StartExecutionInput{
		Input:           cmd.Input,
		Name:            cmd.WorkflowId, // aws step functions use ID as the name (unique identifier for each execution, and can be used as idempotency key)
		StateMachineArn: &cmd.WorkflowName,
	})

	if err != nil {
		fmt.Println(err)
		fmt.Println(output)
		return "", fmt.Errorf("failed to start workflow: %w", err)
	}

	return *output.ExecutionArn, nil
}

type WorkflowStatus string

const (
	Running   WorkflowStatus = "RUNNING"
	Succeeded WorkflowStatus = "SUCCEEDED"
	Failed    WorkflowStatus = "FAILED"
	TimedOut  WorkflowStatus = "TIMED_OUT"
	Cancelled WorkflowStatus = "ABORTED"
)

type WorkflowExecution struct {
	Status WorkflowStatus
	Error  *string
	Output *string
}

func (api *stepFunctionWorkflowApi) GetWorkflowStatus(ctx context.Context, workflowId string) (*WorkflowExecution, error) {
	info, err := api.DescribeExecution(ctx, &sfn.DescribeExecutionInput{
		ExecutionArn: &workflowId,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get workflow status: %w", err)
	}

	return &WorkflowExecution{
		Status: WorkflowStatus(info.Status),
		Error:  info.Error,
		Output: info.Output,
	}, nil
}

func NewWorkflowEngine(ctx context.Context, cfg WorkflowEngineConfig) (WorkflowApi, error) {
	awsCfg, err := newAwsCfg(ctx, cfg.Endpoint, cfg.Region)
	fmt.Println(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build workflow engine: %w", err)
	}

	cli := sfn.NewFromConfig(*awsCfg)
	return &stepFunctionWorkflowApi{cli}, nil
}

type WorkflowEngineConfig struct {
	Region   string
	Endpoint string
}

type stepFunctionWorkflowApi struct {
	*sfn.Client
}

func newAwsCfg(ctx context.Context, endpoint, region string) (*aws.Config, error) {

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			if endpoint != "" {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			}

			// returning EndpointNotFoundError will allow the service to fallback to its default resolution
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)

	}
	return &awsCfg, nil
}
