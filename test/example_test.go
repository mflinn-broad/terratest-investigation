package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/ssh"

	"github.com/gruntwork-io/terratest/modules/retry"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"

	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestTerraformExample(t *testing.T) {
	t.Parallel()
	exampleDir := test_structure.CopyTerraformFolderToTemp(t, "../", "example")

	// get project id
	projectId := gcp.GetGoogleProjectIDFromEnvVar(t)

	// zone for test resources
	zone := "us-east1-b"

	// generate a random name for test bucket
	expectedBucketName := fmt.Sprintf("terratest-gcp-example-%s", strings.ToLower(random.UniqueId()))

	expectedInstanceName := fmt.Sprintf("terratest-gcp-example-%s", strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		//Path to terraform code
		TerraformDir: exampleDir,

		// variables to pass with -var
		Vars: map[string]interface{}{
			"google_project": projectId,
			"zone":           zone,
			"instance_name":  expectedInstanceName,
			"bucket_name":    expectedBucketName,
		},
	}

	// ensure terraform destroy is run
	defer terraform.Destroy(t, terraformOptions)

	// attempt running init and  apply, will fail if errors
	terraform.InitAndApply(t, terraformOptions)

	// get output values
	instanceName := terraform.Output(t, terraformOptions, "instance_name")

	// assert storage bucket exists
	gcp.AssertStorageBucketExists(t, expectedBucketName)

	// Add a tag to instance
	instance := gcp.FetchInstance(t, projectId, instanceName)
	instance.SetLabels(t, map[string]string{"testing": "testing-tag-value2"})

	// check labels with retry
	maxRetries := 12
	timeBetweenRetries := 5 * time.Second
	expectedText := "testing-tag-value2"

	// Check if the instance has the expected tag
	retry.DoWithRetry(t, fmt.Sprintf("Checking Instance %s for labels", instanceName), maxRetries, timeBetweenRetries, func() (string, error) {
		// look up the tags for an instance
		instance := gcp.FetchInstance(t, projectId, instanceName)
		instanceLabels := instance.GetLabels(t)

		testingTag, ok := instanceLabels["testing"]
		actualText := strings.TrimSpace(testingTag)
		if !ok {
			return "", fmt.Errorf("Expected the tag 'testing' to exits")
		}

		if actualText != expectedText {
			return "", fmt.Errorf("Expected GetLabelsForComputeInstanceE to return '%s' but got '%s'", expectedText, actualText)
		}
		return "", nil
	})
}

func TestSSHAccess(t *testing.T) {
	t.Parallel()

	exampleDir := test_structure.CopyTerraformFolderToTemp(t, "../", "example")

	projectID := gcp.GetGoogleProjectIDFromEnvVar(t)
	randomValidGcpName := gcp.RandomValidGcpName()

	// according to docs framework has issue with asia-east2 so exclude it
	zone := gcp.GetRandomZone(t, projectID, nil, nil, []string{"asia-east2"})

	terraformOptions := &terraform.Options{
		TerraformDir: exampleDir,

		// tfvars
		Vars: map[string]interface{}{
			"google_project": projectID,
			"instance_name":  randomValidGcpName,
			"bucket_name":    randomValidGcpName,
			"zone":           zone,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// get pubIp for test instance
	publicIp := terraform.Output(t, terraformOptions, "public_ip")

	instance := gcp.FetchInstance(t, projectID, randomValidGcpName)

	testMessage := "Hello World"
	sshUsername := "terratest"

	keyPair := ssh.GenerateRSAKeyPair(t, 2048)
	instance.AddSshKey(t, sshUsername, keyPair.PublicKey)

	host := ssh.Host{
		Hostname:    publicIp,
		SshKeyPair:  keyPair,
		SshUserName: sshUsername,
	}

	maxRetries := 20
	sleepBetweenRetries := 3 * time.Second

	retry.DoWithRetry(t, "Attempting to SSH", maxRetries, sleepBetweenRetries, func() (string, error) {
		output, err := ssh.CheckSshCommandE(t, host, fmt.Sprintf("echo '%s'", testMessage))
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(testMessage) != strings.TrimSpace(output) {
			return "", fmt.Errorf("Expected: %s. Got: %s\n", testMessage, output)
		}

		return "", nil
	})
}
