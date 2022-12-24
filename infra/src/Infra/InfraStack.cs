using Amazon.CDK;
using Amazon.CDK.AWS.IAM;
using Amazon.CDK.AWS.Lambda;
using Amazon.CDK.AWS.StepFunctions;
using Constructs;
using Amazon.CDK.AWS.StepFunctions.Tasks;

namespace Infra
{

  internal sealed class InfraStack : Stack
  {
    internal InfraStack(Construct scope, string id, IStackProps props = null) : base(scope, id, props)
    {
      var (FindUsersToNotify, NotifyUser) = SetupLambdaFunctions();

      var findUsersToNotify = new LambdaInvoke
      (
        scope: this, 
        id: "find-users-to-notify", 
        props: new LambdaInvokeProps
        {
          LambdaFunction = FindUsersToNotify,
          OutputPath = "$.Payload",
        }
      );

      var notifyUser = new LambdaInvoke
      (
        scope: this, 
        id: "notify-user", 
        props: new LambdaInvokeProps { LambdaFunction = NotifyUser }
      );

      notifyUser.AddRetry
      (
        new RetryProps
        {
          Errors = new []{ "RetryableError"},
          Interval = Duration.Seconds(1),
          MaxAttempts = 3,
          BackoffRate = 2,
        }
      );

      var fanOutUserNotifications = new Map
      (
        scope: this, 
        id: "Fan out user notifications", 
        props: new MapProps
        {
          MaxConcurrency = 5,
          ItemsPath = "$",
          Comment = "Notify users",
        }
      ).Iterator(notifyUser);

    

      var notifyUsersWorkflow = new StateMachine
      (
        scope: this, 
        id: "notifications-workflow",
        props: new StateMachineProps
        {
          StateMachineName = "notifications-workflow",
          Definition =  findUsersToNotify.Next(fanOutUserNotifications)
        }
      );
      
    }


    private (Function FindUsersToNotify, Function NotifyUser) SetupLambdaFunctions()
    {
      var generateLambdaRole = (string name) =>
        new Role
        (
          scope: this, 
          id: $"{name}-role", 
          props: new RoleProps
          {
            AssumedBy = new ServicePrincipal("lambda.amazonaws.com"),
            ManagedPolicies = new[] 
            {
                ManagedPolicy.FromManagedPolicyArn
                (
                  scope: this, 
                  id: $"{name}-role-policy",
                  managedPolicyArn: "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
                )
            }
          }
        );


      var findUsersToNotifyFn = new Function
      (
        scope: this, 
        id: "FindUsersToNotify", 
        props: new FunctionProps
        {
          Runtime = Runtime.NODEJS_14_X,
          Role = generateLambdaRole("find-users-to-notify"),
          Handler = "findUsersToNotify.findUsersToNotify",
          FunctionName = "FindUsersToNotify",
          Code = Code.FromAsset("../src/dist/findUsersToNotify.zip")
        }
      );

      var notifyUserFn = new Function(
        scope: this, 
        id: "NotifyUser", 
        props: new FunctionProps
        {
          Runtime = Runtime.NODEJS_14_X,
          Role = generateLambdaRole("notify-user"),
          Handler = "notifyUser.notifyUser",
          FunctionName = "NotifyUser",
          Code = Code.FromAsset("../src/dist/notifyUser.zip")
        }
      );
      return (findUsersToNotifyFn, notifyUserFn);
    }
  }
}

