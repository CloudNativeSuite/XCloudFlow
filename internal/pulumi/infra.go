package pulumi

import (
	"fmt"

	awsVPC "PulumiGo/internal/pulumi/modules/aws/vpc"
	"PulumiGo/internal/pulumi/utils"
	pulumisdk "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// DeployInfrastructure loads configuration and provisions resources.
func DeployInfrastructure() error {
	cfg, err := utils.LoadMergedConfig("")
	if err != nil {
		return err
	}

	return pulumisdk.RunErr(func(ctx *pulumisdk.Context) error {
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
	})
}
