package main

import (
	"sync"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

type mocks int

func (mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return args.Name + "_id", args.Inputs, nil
}

func (mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}

func TestInfrastructure(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		infra, err := createInfrastructure(ctx)
		assert.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(3)

		// test if create vpc with cidr correctly
		pulumi.All(infra.Vpc.URN(), infra.Vpc.CidrBlock).ApplyT(func(all []interface{}) error {
			urn := all[0].(pulumi.URN)
			cidr := all[1].(string)
			assert.Equal(t, cidr, "10.2.0.0/16", "missing or wrong cidr: %s", urn)
			wg.Done()
			return nil
		})

		//test if create subnet under vpc correctly
		pulumi.All(infra.Vpc.URN(), infra.Vpc.ID(), infra.Subnets[0].VpcId, infra.Subnets[1].VpcId).ApplyT(func(all []interface{}) error {
			urn := all[0].(pulumi.URN)
			vpcId := string(all[1].(pulumi.ID))
			subnetVpcId1 := all[2].(string)
			subnetVpcId2 := all[3].(string)

			assert.Equal(t, subnetVpcId1, vpcId, "missing or wrong vpc: %s", urn)
			assert.Equal(t, subnetVpcId2, vpcId, "missing or wrong vpc: %s", urn)
			wg.Done()
			return nil
		})

		// test if create subnet with cidr correctly
		pulumi.All(infra.Vpc.URN(), infra.Subnets[0].CidrBlock, infra.Subnets[1].CidrBlock).ApplyT(func(all []interface{}) error {
			urn := all[0].(pulumi.URN)
			subnetCidrBlock1 := all[1].(string)
			subnetCidrBlock2 := all[2].(string)

			assert.Equal(t, subnetCidrBlock1, "10.2.0.0/24", "missing or wrong cidr: %s", urn)
			assert.Equal(t, subnetCidrBlock2, "10.2.1.0/24", "missing or wrong cidr: %s", urn)
			wg.Done()
			return nil
		})

		wg.Wait()
		return nil
	}, pulumi.WithMocks("avantis-aws-infrastructure", "dev", mocks(0)))
	assert.NoError(t, err)
}
