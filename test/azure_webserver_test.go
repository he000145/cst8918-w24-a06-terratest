package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "866e3e98-5089-4944-8cd1-2f0c2dfddad0"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "he000145",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variables
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Confirm NIC exists
	assert.True(t, azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID))

	// Confirm the NIC is attached to the VM
	attachedNics := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)
	assert.Contains(t, attachedNics, nicName)

	// Confirm VM is using the expected Ubuntu image
	vmImage := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)
	assert.Equal(t, "Canonical", vmImage.Publisher)
	assert.Equal(t, "0001-com-ubuntu-server-jammy", vmImage.Offer)
	assert.Equal(t, "22_04-lts-gen2", vmImage.SKU)
}
