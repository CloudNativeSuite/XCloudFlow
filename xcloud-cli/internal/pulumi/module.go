package pulumi

import (
	"context"
	"encoding/json"

	"xcloud-cli/internal/modules"
)

// DeployTask triggers Pulumi infrastructure deployment.
type DeployTask struct{}

// Type implements modules.Task
func (DeployTask) Type() string { return "pulumi_deploy" }

// deployHandler runs DeployInfrastructure and records the result.
type deployHandler struct{}

func (deployHandler) Run(ctx context.Context, t modules.Task) (string, error) {
	res, err := DeployInfrastructure()
	if err != nil {
		return "", err
	}
	for _, stack := range res.Targets {
		if stack.Status == "applied" {
			modules.RecordResource(ctx, modules.Resource{ID: stack.Stack, Type: "pulumi", Name: stack.Cloud + "/" + stack.Region})
		}
	}

	payload, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func init() {
	modules.Register("pulumi_deploy", deployHandler{})
}
