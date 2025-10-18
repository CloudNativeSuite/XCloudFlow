package vpc

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	pulumi "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SubnetConfig struct {
	Name             string `yaml:"name"`
	CidrBlock        string `yaml:"cidr_block"`
	Type             string `yaml:"type"`
	AvailabilityZone string `yaml:"availability_zone"`
}

type RouteConfig struct {
	DestinationCidrBlock string `yaml:"destination_cidr_block"`
	SubnetType           string `yaml:"subnet_type"`
}

type VpcConfig struct {
	Name      string         `yaml:"name"`
	CidrBlock string         `yaml:"cidr_block"`
	Subnets   []SubnetConfig `yaml:"subnets"`
	Routes    []RouteConfig  `yaml:"routes"`
}

type VpcResult struct {
	Vpc     *ec2.Vpc
	Subnets map[string]*ec2.Subnet
	Igw     *ec2.InternetGateway
}

func CreateVPC(ctx *pulumi.Context, conf VpcConfig, region string) (*VpcResult, error) {
	vpc, err := ec2.NewVpc(ctx, conf.Name, &ec2.VpcArgs{
		CidrBlock: pulumi.String(conf.CidrBlock),
		Tags:      pulumi.StringMap{"Name": pulumi.String(conf.Name)},
	})
	if err != nil {
		return nil, err
	}

	hasPublic := false
	for _, s := range conf.Subnets {
		if s.Type == "public" {
			hasPublic = true
			break
		}
	}

	var igw *ec2.InternetGateway
	if hasPublic {
		igw, err = ec2.NewInternetGateway(ctx, "main-igw", &ec2.InternetGatewayArgs{
			VpcId: vpc.ID(),
		})
		if err != nil {
			return nil, err
		}
	}

	subnets := make(map[string]*ec2.Subnet)
	for _, s := range conf.Subnets {
		subnet, err := ec2.NewSubnet(ctx, s.Name, &ec2.SubnetArgs{
			VpcId:               vpc.ID(),
			CidrBlock:           pulumi.String(s.CidrBlock),
			MapPublicIpOnLaunch: pulumi.Bool(s.Type == "public"),
			AvailabilityZone:    pulumi.String(s.AvailabilityZone),
			Tags:                pulumi.StringMap{"Name": pulumi.String(s.Name)},
		})
		if err != nil {
			return nil, err
		}
		subnets[s.Name] = subnet
	}

	if hasPublic {
		var routes ec2.RouteTableRouteArray
		for _, r := range conf.Routes {
			if r.SubnetType == "public" {
				routes = append(routes, ec2.RouteTableRouteArgs{
					CidrBlock: pulumi.String(r.DestinationCidrBlock),
					GatewayId: igw.ID(),
				})
			}
		}
		rt, err := ec2.NewRouteTable(ctx, "public-route-table", &ec2.RouteTableArgs{
			VpcId:  vpc.ID(),
			Routes: routes,
		})
		if err != nil {
			return nil, err
		}

		for _, s := range conf.Subnets {
			if s.Type == "public" {
				_, err = ec2.NewRouteTableAssociation(ctx, s.Name+"-assoc", &ec2.RouteTableAssociationArgs{
					SubnetId:     subnets[s.Name].ID(),
					RouteTableId: rt.ID(),
				})
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return &VpcResult{Vpc: vpc, Subnets: subnets, Igw: igw}, nil
}
