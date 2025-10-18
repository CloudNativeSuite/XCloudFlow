package pulumi

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	auto "github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	pulumiSdk "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"xcloud-cli/internal/modules/utils"
	awsVPC "xcloud-cli/internal/pulumi/modules/aws/vpc"
)

// DeploymentResult captures aggregated information from a Pulumi run.
type DeploymentResult struct {
	ActiveEnv     string            `json:"activeEnv"`
	ConfigSources []string          `json:"configSources"`
	Targets       []StackDeployment `json:"targets"`
}

// StackDeployment represents the outcome for a single cloud/region stack.
type StackDeployment struct {
	Stack        string            `json:"stack"`
	Cloud        string            `json:"cloud"`
	Region       string            `json:"region"`
	Status       string            `json:"status"`
	Message      string            `json:"message,omitempty"`
	PlanPreview  map[string]int    `json:"planPreview,omitempty"`
	Applied      map[string]int    `json:"applied,omitempty"`
	CostEstimate *CostEstimate     `json:"costEstimate,omitempty"`
	Outputs      map[string]string `json:"outputs,omitempty"`
	Artifacts    []string          `json:"artifacts,omitempty"`
	PreviewURL   string            `json:"previewUrl,omitempty"`
	UpdateURL    string            `json:"updateUrl,omitempty"`
}

// CostEstimate provides a lightweight view of expected spend.
type CostEstimate struct {
	HourlyUSD   float64            `json:"hourlyUSD"`
	MonthlyUSD  float64            `json:"monthlyUSD"`
	Breakdown   map[string]float64 `json:"breakdownUSD"`
	Assumptions []string           `json:"assumptions"`
}

// DeploymentSpec captures the matrix definition from the DSL.
type DeploymentSpec struct {
	Org           string       `yaml:"org" json:"org"`
	Project       string       `yaml:"project" json:"project"`
	DefaultRegion string       `yaml:"defaultRegion" json:"defaultRegion"`
	Matrix        MatrixConfig `yaml:"matrix" json:"matrix"`
}

// MatrixConfig defines the clouds/regions expansion.
type MatrixConfig struct {
	Clouds  []string            `yaml:"clouds" json:"clouds"`
	Regions map[string][]string `yaml:"regions" json:"regions"`
}

type deploymentTarget struct {
	Cloud  string
	Region string
}

// infraProgram returns the Pulumi program used for deployments.
func infraProgram(cfg map[string]interface{}, target deploymentTarget) pulumiSdk.RunFunc {
	return func(ctx *pulumiSdk.Context) error {
		ctx.Export("cloud", pulumiSdk.String(target.Cloud))
		ctx.Export("region", pulumiSdk.String(target.Region))

		if target.Cloud != "aws" {
			ctx.Log.Info(fmt.Sprintf("Skipping unsupported cloud provider %s", target.Cloud), nil)
			return nil
		}

		var vpcConf awsVPC.VpcConfig
		if err := utils.DecodeSection(cfg, "vpc", &vpcConf); err != nil {
			return err
		}

		res, err := awsVPC.CreateVPC(ctx, vpcConf, target.Region)
		if err != nil {
			return err
		}
		ctx.Export("vpcId", res.Vpc.ID())
		fmt.Printf("✅ Pulumi stack %s deployed for %s/%s.\n", ctx.Stack(), target.Cloud, target.Region)
		return nil
	}
}

// DeployInfrastructure provisions resources using the Pulumi Automation API.
func DeployInfrastructure() (*DeploymentResult, error) {
	ctx := context.Background()
	env := os.Getenv("STACK_ENV")
	if env == "" {
		env = "sit"
	}

	cfg, err := utils.LoadMergedConfig("")
	if err != nil {
		return nil, err
	}

	spec, err := extractDeploymentSpec(cfg)
	if err != nil {
		return nil, err
	}

	targets, err := expandTargets(spec, os.Getenv("STACK_CLOUD"), os.Getenv("STACK_REGION"))
	if err != nil {
		return nil, err
	}

	configSources := []string{}
	if raw, ok := cfg["__config_path__"].(string); ok && raw != "" {
		for _, part := range strings.Split(raw, ",") {
			configSources = append(configSources, part)
		}
	}

	deployments := make([]StackDeployment, 0, len(targets))
	for _, target := range targets {
		if target.Cloud != "aws" {
			deployments = append(deployments, StackDeployment{
				Stack:   buildStackName(env, target),
				Cloud:   target.Cloud,
				Region:  target.Region,
				Status:  "skipped",
				Message: fmt.Sprintf("cloud provider %s not yet supported", target.Cloud),
			})
			continue
		}

		summary, err := runStack(ctx, env, spec.Project, cfg, target)
		if err != nil {
			deployments = append(deployments, StackDeployment{
				Stack:   buildStackName(env, target),
				Cloud:   target.Cloud,
				Region:  target.Region,
				Status:  "failed",
				Message: err.Error(),
			})
			continue
		}
		deployments = append(deployments, *summary)
	}

	return &DeploymentResult{
		ActiveEnv:     env,
		ConfigSources: configSources,
		Targets:       deployments,
	}, nil
}

func buildStackName(env string, target deploymentTarget) string {
	sanitizedRegion := strings.NewReplacer("/", "-", " ", "-", "_", "-").Replace(target.Region)
	return fmt.Sprintf("%s-%s-%s", env, target.Cloud, sanitizedRegion)
}

func runStack(ctx context.Context, env, project string, cfg map[string]interface{}, target deploymentTarget) (*StackDeployment, error) {
	projectName := project
	if projectName == "" {
		projectName = "xcloud"
	}

	stackName := buildStackName(env, target)
	stack, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, infraProgram(cfg, target))
	if err != nil {
		return nil, err
	}

	if err := stack.Workspace().SetEnvVars(map[string]string{
		"PULUMI_CONFIG_PASSPHRASE_FILE": filepath.Join(os.Getenv("HOME"), ".pulumi-passphrase"),
		"CONFIG_PATH":                   os.Getenv("CONFIG_PATH"),
		"STACK_ENV":                     env,
		"STACK_CLOUD":                   target.Cloud,
		"STACK_REGION":                  target.Region,
	}); err != nil {
		return nil, err
	}

	previewRes, err := stack.Preview(ctx)
	if err != nil {
		return nil, err
	}

	upRes, err := stack.Up(ctx)
	if err != nil {
		return nil, err
	}

	plan := summarizeChanges(previewRes.ChangeSummary)
	applied := map[string]int{}
	if upRes.Summary.ResourceChanges != nil {
		applied = copyStringIntMap(*upRes.Summary.ResourceChanges)
	}

	outputs := make(map[string]string)
	keys := make([]string, 0, len(upRes.Outputs))
	for k := range upRes.Outputs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		outputs[k] = fmt.Sprint(upRes.Outputs[k].Value)
	}

	previewURL, _ := previewRes.GetPermalink()
	updateURL, _ := upRes.GetPermalink()

	artifacts := []string{
		fmt.Sprintf("Stack state file: %s", filepath.Join(stack.Workspace().WorkDir(), fmt.Sprintf("Pulumi.%s.yaml", stackName))),
	}

	cost := estimateCost(cfg)
	return &StackDeployment{
		Stack:        stackName,
		Cloud:        target.Cloud,
		Region:       target.Region,
		Status:       "applied",
		PlanPreview:  plan,
		Applied:      applied,
		CostEstimate: &cost,
		Outputs:      outputs,
		Artifacts:    artifacts,
		PreviewURL:   previewURL,
		UpdateURL:    updateURL,
	}, nil
}

func summarizeChanges(summary map[apitype.OpType]int) map[string]int {
	out := make(map[string]int, len(summary))
	for k, v := range summary {
		out[string(k)] = v
	}
	return out
}

func copyStringIntMap(m map[string]int) map[string]int {
	out := make(map[string]int, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func extractDeploymentSpec(cfg map[string]interface{}) (*DeploymentSpec, error) {
	spec := &DeploymentSpec{}
	if v, ok := cfg["org"].(string); ok {
		spec.Org = v
	}
	if v, ok := cfg["project"].(string); ok {
		spec.Project = v
	}
	if v, ok := cfg["defaultRegion"].(string); ok {
		spec.DefaultRegion = v
	}
	if spec.DefaultRegion == "" {
		if awsSection, ok := cfg["aws"].(map[string]interface{}); ok {
			if region, ok := awsSection["region"].(string); ok {
				spec.DefaultRegion = region
			}
		}
	}

	if err := utils.DecodeSection(cfg, "matrix", &spec.Matrix); err != nil {
		return nil, fmt.Errorf("matrix configuration missing or invalid: %w", err)
	}

	if len(spec.Matrix.Clouds) == 0 {
		for cloud := range spec.Matrix.Regions {
			spec.Matrix.Clouds = append(spec.Matrix.Clouds, cloud)
		}
		sort.Strings(spec.Matrix.Clouds)
	}

	if spec.DefaultRegion == "" {
		return nil, errors.New("default region is required when matrix regions omit entries")
	}
	return spec, nil
}

func expandTargets(spec *DeploymentSpec, overrideCloud, overrideRegion string) ([]deploymentTarget, error) {
	if overrideCloud != "" {
		region := overrideRegion
		if region == "" {
			if regions, ok := spec.Matrix.Regions[overrideCloud]; ok && len(regions) > 0 {
				region = regions[0]
			} else {
				region = spec.DefaultRegion
			}
		}
		return []deploymentTarget{{Cloud: overrideCloud, Region: region}}, nil
	}

	if len(spec.Matrix.Clouds) == 0 {
		return nil, errors.New("matrix configuration requires at least one cloud")
	}

	var targets []deploymentTarget
	for _, cloud := range spec.Matrix.Clouds {
		regions := spec.Matrix.Regions[cloud]
		if len(regions) == 0 {
			regions = []string{spec.DefaultRegion}
		}
		if overrideRegion != "" {
			regions = []string{overrideRegion}
		}
		for _, region := range regions {
			targets = append(targets, deploymentTarget{Cloud: cloud, Region: region})
		}
	}
	return targets, nil
}

func estimateCost(cfg map[string]interface{}) CostEstimate {
	instances, _ := cfg["instances"].([]interface{})
	breakdown := make(map[string]float64)
	totalHourly := 0.0
	for _, inst := range instances {
		data, ok := inst.(map[string]interface{})
		if !ok {
			continue
		}
		instanceType, _ := data["type"].(string)
		lifecycle, _ := data["lifecycle"].(string)
		hourly := lookupInstancePrice(instanceType)
		if lifecycle == "spot" {
			hourly *= 0.35
		}
		totalHourly += hourly
		breakdown[instanceType] += hourly
	}

	monthly := totalHourly * 730
	for k, v := range breakdown {
		breakdown[k] = math.Round(v*730*100) / 100
	}

	assumptions := []string{
		"730 hours per month",
		"Spot instances estimated at 35% of on-demand pricing",
	}
	return CostEstimate{
		HourlyUSD:   math.Round(totalHourly*100) / 100,
		MonthlyUSD:  math.Round(monthly*100) / 100,
		Breakdown:   breakdown,
		Assumptions: assumptions,
	}
}

func lookupInstancePrice(instanceType string) float64 {
	if instanceType == "" {
		return 0
	}
	prices := map[string]float64{
		"t3.micro":  0.0104,
		"t3.small":  0.0208,
		"t3.medium": 0.0416,
		"t3.large":  0.0832,
	}
	if p, ok := prices[instanceType]; ok {
		return p
	}
	return 0.02
}
