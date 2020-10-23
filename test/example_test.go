package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/retry"

	"github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terratest/modules/terraform"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"

	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestTerraformExample(t *testing.T) {

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
	bucketURL := terraform.Output(t, terraformOptions, "bucket_url")
	instanceName := terraform.Output(t, terraformOptions, "instance_name")

	expectedURL := fmt.Sprintf("gs://%s", expectedBucketName)
	assert.Equal(t, expectedURL, bucketURL)

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
