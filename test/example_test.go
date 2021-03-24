package test

import (
	"fmt"
	"strings"
	"testing"

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
}
