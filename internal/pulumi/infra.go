package pulumi

import (
	"context"
	"fmt"
	"os"

	"PulumiGo/internal/modules/utils"
	awsVPC "PulumiGo/internal/pulumi/modules/aws/vpc"
	auto "github.com/pulumi/pulumi/sdk/v3/go/auto"
	pulumiSdk "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// infraProgram returns the Pulumi program used for deployments.
func infraProgram() pulumiSdk.RunFunc {
	return func(ctx *pulumiSdk.Context) error {
		cfg, err := utils.LoadMergedConfig("")
		if err != nil {
			return err
		}

		var vpcConf awsVPC.VpcConfig
		if err := utils.DecodeSection(cfg, "vpc", &vpcConf); err != nil {
			return err
		}

		res, err := awsVPC.CreateVPC(ctx, vpcConf, ctx.Stack())
		if err != nil {
			return err
		}
		ctx.Export("vpcId", res.Vpc.ID())
		fmt.Println("✅ Pulumi stack deployed.")
		return nil
	}
}

// DeployInfrastructure provisions resources using the Pulumi Automation API.
func DeployInfrastructure() error {
	ctx := context.Background()
	env := os.Getenv("STACK_ENV")
	if env == "" {
		env = "sit"
	}

	stack, err := auto.UpsertStackInlineSource(ctx, env, "PulumiGo", infraProgram())
	if err != nil {
		return err
	}

	stack.Workspace().SetEnvVars(map[string]string{
		"PULUMI_CONFIG_PASSPHRASE_FILE": os.Getenv("HOME") + "/.pulumi-passphrase",
		"CONFIG_PATH":                   os.Getenv("CONFIG_PATH"),
		"STACK_ENV":                     env,
	})

	_, err = stack.Up(ctx)
	return err
}
