package securitygroup

import (
	"strconv"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	pulumi "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AllowRule struct {
	Protocol string        `yaml:"protocol"`
	Ports    []interface{} `yaml:"ports"`
}

type GroupRuleConfig struct {
	Name         string      `yaml:"name"`
	SourceRanges []string    `yaml:"source_ranges"`
	EgressRanges []string    `yaml:"egress_ranges"`
	Allow        []AllowRule `yaml:"allow"`
}

// CreateSecurityGroup creates a security group with ingress and egress rules.
func CreateSecurityGroup(ctx *pulumi.Context, vpcID pulumi.StringInput, cfg GroupRuleConfig) (*ec2.SecurityGroup, error) {
	ingress := ec2.SecurityGroupIngressArray{}

	srcRanges := cfg.SourceRanges
	if len(srcRanges) == 0 {
		srcRanges = []string{"0.0.0.0/0"}
	}

	for _, rule := range cfg.Allow {
		proto := rule.Protocol
		if proto == "" {
			proto = "tcp"
		}
		for _, p := range rule.Ports {
			var from, to int
			switch val := p.(type) {
			case int:
				from, to = val, val
			case int64:
				from, to = int(val), int(val)
			case float64:
				from, to = int(val), int(val)
			case string:
				lp := strings.ToLower(val)
				if lp == "*" || lp == "any" || lp == "all" {
					from, to = 0, 65535
				} else {
					if n, err := strconv.Atoi(val); err == nil {
						from, to = n, n
					} else {
						continue
					}
				}
			default:
				continue
			}
			ingress = append(ingress, ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String(proto),
				FromPort:   pulumi.Int(from),
				ToPort:     pulumi.Int(to),
				CidrBlocks: pulumi.ToStringArray(srcRanges),
			})
		}
	}

	egressRanges := cfg.EgressRanges
	if len(egressRanges) == 0 {
		egressRanges = []string{"0.0.0.0/0"}
	}

	sg, err := ec2.NewSecurityGroup(ctx, cfg.Name, &ec2.SecurityGroupArgs{
		VpcId:   vpcID,
		Ingress: ingress,
		Egress: ec2.SecurityGroupEgressArray{
			ec2.SecurityGroupEgressArgs{
				Protocol:   pulumi.String("-1"),
				FromPort:   pulumi.Int(0),
				ToPort:     pulumi.Int(0),
				CidrBlocks: pulumi.ToStringArray(egressRanges),
			},
		},
		Tags: pulumi.StringMap{
			"Name": pulumi.String(cfg.Name),
		},
	})
	if err != nil {
		return nil, err
	}
	return sg, nil
}
