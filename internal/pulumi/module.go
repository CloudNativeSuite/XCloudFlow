package pulumi

import (
	"context"

	"PulumiGo/internal/modules"
)

// DeployTask triggers Pulumi infrastructure deployment.
type DeployTask struct{}

// Type implements modules.Task
func (DeployTask) Type() string { return "pulumi_deploy" }

// deployHandler runs DeployInfrastructure and records the result.
type deployHandler struct{}

func (deployHandler) Run(ctx context.Context, t modules.Task) (string, error) {
	if err := DeployInfrastructure(); err != nil {
		return "", err
	}
	modules.RecordResource(ctx, modules.Resource{ID: "stack", Type: "pulumi", Name: "demo"})
	return "Pulumi stack deployed", nil
}

func init() {
	modules.Register("pulumi_deploy", deployHandler{})
}
