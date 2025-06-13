package ec2module

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	pulumi "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type InstanceConfig struct {
	Name            string `yaml:"name"`
	Subnet          string `yaml:"subnet"`
	Ami             string `yaml:"ami"`
	Type            string `yaml:"type"`
	DiskSizeGb      int    `yaml:"disk_size_gb"`
	Lifecycle       string `yaml:"lifecycle"`
	Ttl             string `yaml:"ttl"`
	Env             string `yaml:"env"`
	Owner           string `yaml:"owner"`
	UserData        string `yaml:"user_data"`
	PrivateIp       string `yaml:"private_ip"`
	AssociatePublic bool   `yaml:"associate_public_ip"`
}

type InstanceOutputs map[string]pulumi.StringOutput

func resolveAMI(ctx *pulumi.Context, keyword string, region string) (string, error) {
	if strings.HasPrefix(keyword, "ami-") {
		return keyword, nil
	}
	kw := strings.ToLower(keyword)
	switch kw {
	case "ubuntu-22.04", "ubuntu22.04":
		res, err := ec2.LookupAmi(ctx, &ec2.LookupAmiArgs{
			Owners:     []string{"099720109477"},
			MostRecent: pulumi.BoolRef(true),
			Filters: []ec2.GetAmiFilter{
				{Name: "name", Values: []string{"ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"}},
				{Name: "virtualization-type", Values: []string{"hvm"}},
			},
		}, nil)
		if err != nil {
			return "", err
		}
		return res.Id, nil
	case "rocky-8.10", "rockylinux-8.10", "rocky8.10":
		res, err := ec2.LookupAmi(ctx, &ec2.LookupAmiArgs{
			Owners:     []string{"792107900819"},
			MostRecent: pulumi.BoolRef(true),
			Filters: []ec2.GetAmiFilter{
				{Name: "name", Values: []string{"Rocky-8-ec2-8.10*x86_64"}},
				{Name: "architecture", Values: []string{"x86_64"}},
			},
		}, nil)
		if err != nil {
			return "", err
		}
		return res.Id, nil
	default:
		return "", fmt.Errorf("unsupported AMI keyword: %s", keyword)
	}
}

func loadUserData(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	fp := filepath.Clean(os.ExpandEnv(path))
	data, err := os.ReadFile(fp)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func CreateInstances(ctx *pulumi.Context, configs []InstanceConfig, subnets map[string]*ec2.Subnet, sg *ec2.SecurityGroup, keyName string, deps []pulumi.Resource) (InstanceOutputs, error) {
	outs := InstanceOutputs{}
	for _, c := range configs {
		subnet, ok := subnets[c.Subnet]
		if !ok {
			return nil, fmt.Errorf("subnet %s not found", c.Subnet)
		}
		ami, err := resolveAMI(ctx, c.Ami, ctx.Region())
		if err != nil {
			return nil, err
		}
		userData, err := loadUserData(c.UserData)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		depsAll := []pulumi.Resource{subnet}
		if sg != nil {
			depsAll = append(depsAll, sg)
		}
		depsAll = append(depsAll, deps...)
		inst, err := ec2.NewInstance(ctx, c.Name, &ec2.InstanceArgs{
			Ami:                      pulumi.StringPtr(ami),
			InstanceType:             pulumi.String(c.Type),
			KeyName:                  pulumi.StringPtr(keyName),
			SubnetId:                 subnet.ID(),
			PrivateIp:                pulumi.StringPtr(c.PrivateIp),
			AssociatePublicIpAddress: pulumi.BoolPtr(c.AssociatePublic),
			VpcSecurityGroupIds:      pulumi.StringArray{sg.ID()},
			UserData:                 pulumi.StringPtr(userData),
			RootBlockDevice: &ec2.InstanceRootBlockDeviceArgs{
				VolumeSize: pulumi.IntPtr(c.DiskSizeGb),
				VolumeType: pulumi.StringPtr("gp2"),
			},
			Tags: pulumi.StringMap{
				"Name":        pulumi.String(c.Name),
				"Lifecycle":   pulumi.String(c.Lifecycle),
				"TTL":         pulumi.String(c.Ttl),
				"Environment": pulumi.String(c.Env),
				"Owner":       pulumi.String(c.Owner),
			},
		}, pulumi.DependsOn(depsAll))
		if err != nil {
			return nil, err
		}
		outs[c.Name+"_id"] = inst.ID().ToStringOutput()
		outs[c.Name+"_public_ip"] = inst.PublicIp
		outs[c.Name+"_private_ip"] = inst.PrivateIp
	}
	return outs, nil
}
