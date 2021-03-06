# terratest-investigation
The purpose of this repo is to demonstrate a basic workflow for automated infrastructure
testing using [terratest](https://terratest.gruntwork.io/)

In this example workflow tests are run against a simple terraform module but
automated testing of other types of infrastructure are possible too.

All of the following are supported in the automated testing framework:
- Testing Terraform code
- Testing Packer templates
- Testing Docker images
- Executing commands on servers over SSH
- Working with AWS APIs
- Working with Azure APIs
- Working with GCP APIs
- Working with Kubernetes APIs
- Testing Helm Charts
- Making HTTP requests
- Running shell commands

## Demo Workflow

### Examples
- [Example successful workflow](https://github.com/mflinn-broad/terratest-investigation/pull/3)
- [Example failed workflow](https://github.com/mflinn-broad/terratest-investigation/pull/4)

This repo incorporates a basic github actions workflow which will run the automated infrastructure
tests on prs to `master` and a passing test is required before merging.

The automated tests are currently configured to use `dsp-tooks-k8s` as the target project for spinning
up test infrastructure. using a minimal service account needed for the example tf module.

The automated test performs a "clean apply" of the example module on each run meaning it starts from
an empty tfstate. This helps to catch interdependencies preventing clean applies that can occur when building up a module
with iterative applies.

The tests will spin up real infrastructure and make assertions about it rather than just checking terraform state.
For example this workflow will test the ability to ssh to GCE VM created by the example module and run a basic shell command

Go's `defer` functionality is used to ensure that even in the case of unexpected failures, all infrastructure that is created
by the tests will be cleaned up, or in other words `terraform destroy` will always run at the end of each test run

To see an example of a successful test run, open a no-op pr to master, the test will automatically run when the pr is opened.
and it must pass before you are allowed to merge

To see an example of a failed test run open a pr where you hardcode the name parameter on either the GCE VM or GCS bucket rather
using a variable

The test currently takes about 5 minutes to run so there may be a bit of a wait
