package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Subnet struct {
	SubnetName string
	SubnetCidr string
}
type SubnetBlocks []Subnet

// Name
const DEV_VPC = "dev-vpc"
const PRIVATE_SUBNET = "private-subnet"
const PUBLIC_SUBNET = "public-subnet"

// CIDR
const DEV_VPC_CIDR = "10.2.0.0/16"
const PRIVATE_SUBNET_CIDR = "10.2.0.0/24"
const PUBLIC_SUBNET_CIDR = "10.2.1.0/24"

type Infrastructure struct {
	Vpc     *ec2.Vpc
	Subnets []*ec2.Subnet
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := createInfrastructure(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

func createInfrastructure(ctx *pulumi.Context) (*Infrastructure, error) {
	devPvc, err := createVpc(ctx, DEV_VPC, DEV_VPC_CIDR)
	if err != nil {
		return nil, err
	}

	subnetsCidrBlocksList := SubnetBlocks{
		Subnet{SubnetName: PRIVATE_SUBNET, SubnetCidr: PRIVATE_SUBNET_CIDR},
		Subnet{SubnetName: PUBLIC_SUBNET, SubnetCidr: PUBLIC_SUBNET_CIDR},
	}

	subnets, err := createListOfSubnetsFromGivenVpc(ctx, devPvc, subnetsCidrBlocksList)
	if err != nil {
		return nil, err
	}

	return &Infrastructure{
		Vpc:     devPvc,
		Subnets: subnets,
	}, nil
}

func createVpc(ctx *pulumi.Context, vpcName string, vpcCidr string) (*ec2.Vpc, error) {
	devVpc, err := ec2.NewVpc(ctx, vpcName, &ec2.VpcArgs{
		CidrBlock: pulumi.String(vpcCidr),
	})

	if err != nil {
		return nil, err
	}

	return devVpc, nil
}

func createListOfSubnetsFromGivenVpc(ctx *pulumi.Context, vpc *ec2.Vpc, listOfSubnets SubnetBlocks) ([]*ec2.Subnet, error) {
	var subnets []*ec2.Subnet
	for _, subnetsCidrBlock := range listOfSubnets {
		var subnet *ec2.Subnet

		subnet, err := ec2.NewSubnet(ctx, subnetsCidrBlock.SubnetName, &ec2.SubnetArgs{
			VpcId:     vpc.ID(),
			CidrBlock: pulumi.String(subnetsCidrBlock.SubnetCidr),
		})

		if err != nil {
			return nil, err
		}
		subnets = append(subnets, subnet)
	}

	return subnets, nil
}
