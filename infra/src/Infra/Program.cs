using Amazon.CDK;
using Infra;

var app = new App();
var _ = new InfraStack(app, "InfraStack");
app.Synth();
