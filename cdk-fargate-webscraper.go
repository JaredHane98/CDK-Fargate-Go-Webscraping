package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapplicationautoscaling"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	ecs_patterns "github.com/aws/aws-cdk-go/awscdk/v2/awsecspatterns"
	"github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkFargateWebscraperStackProps struct {
	awscdk.StackProps
}

func NewCdkFargateWebscraperStack(scope constructs.Construct, id string, props *CdkFargateWebscraperStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	dockerImageAsset := awsecrassets.NewDockerImageAsset(stack, jsii.String("MyDockerImage"), &awsecrassets.DockerImageAssetProps{
		Directory: jsii.String("./Container/"),
	})

	vpc := awsec2.NewVpc(stack, jsii.String("ALBFargoVpc"), &awsec2.VpcProps{
		MaxAzs:      jsii.Number(2),
		IpAddresses: awsec2.IpAddresses_Cidr(jsii.String("10.10.0.0/16")),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				SubnetType: awsec2.SubnetType_PUBLIC,
				Name:       jsii.String("Public"),
				CidrMask:   jsii.Number(24),
			},
			{
				SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
				Name:       jsii.String("Private"),
				CidrMask:   jsii.Number(24),
			},
		},
		NatGateways: jsii.Number(1), // using a nat gateway to provide outbound access
	})

	cluster := awsecs.NewCluster(stack, jsii.String("ALBFargoECSCluster"), &awsecs.ClusterProps{
		Vpc: vpc,
	})

	fargateService := ecs_patterns.NewApplicationLoadBalancedFargateService(stack, jsii.String("ALBFargoService"), &ecs_patterns.ApplicationLoadBalancedFargateServiceProps{
		Cluster: cluster,
		TaskImageOptions: &ecs_patterns.ApplicationLoadBalancedTaskImageOptions{
			Image: awsecs.ContainerImage_FromDockerImageAsset(dockerImageAsset),
		},
		RuntimePlatform: &awsecs.RuntimePlatform{
			OperatingSystemFamily: awsecs.OperatingSystemFamily_LINUX(),
			CpuArchitecture:       awsecs.CpuArchitecture_X86_64(),
		},
		DesiredCount:       jsii.Number(1),
		Cpu:                jsii.Number(256),
		MemoryLimitMiB:     jsii.Number(512),
		PublicLoadBalancer: jsii.Bool(true),
		ListenerPort:       jsii.Number(80),

		TaskSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
		},
	})

	scaling := fargateService.Service().AutoScaleTaskCount(&awsapplicationautoscaling.EnableScalingProps{
		MinCapacity: jsii.Number(1),
		MaxCapacity: jsii.Number(5),
	})

	scaling.ScaleOnCpuUtilization(jsii.String("CpuScaling"), &awsecs.CpuUtilizationScalingProps{
		TargetUtilizationPercent: jsii.Number(50),
	})
	scaling.ScaleOnMemoryUtilization(jsii.String("MemoryScaling"), &awsecs.MemoryUtilizationScalingProps{
		TargetUtilizationPercent: jsii.Number(50),
	})

	fargateService.TargetGroup().ConfigureHealthCheck(&awselasticloadbalancingv2.HealthCheck{
		Path: jsii.String("/health"),
	})

	awscdk.NewCfnOutput(stack, jsii.String("LoadBalancerDNS"), &awscdk.CfnOutputProps{Value: fargateService.LoadBalancer().LoadBalancerDnsName()})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkFargateWebscraperStack(app, "CdkFargateWebscraperStack", &CdkFargateWebscraperStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
