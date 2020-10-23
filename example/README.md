# Terratest Demo Workflow module

This is a simple terra form module which create a GCE VM instance and GCS bucket.  
It is intended to be used just for a demonstration of a sample automated infrasture  
testing workflow using [terratest](https://terratest.gruntwork.io/)

To break the automated tests make a pr with the instance or bucket name hard coded  
rather than using a variable

This documentation is generated with [terraform-docs](https://github.com/segmentio/terraform-docs)
`terraform-docs markdown --no-sort . > README.md`

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 0.12.26 |

## Providers

| Name | Version |
|------|---------|
| google | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| google\_project | google project to run tests in | `string` | n/a | yes |
| instance\_name | name of gcp vm instance | `string` | `"terratest-example"` | no |
| machine\_type | Machine type of vm | `string` | `"f1-micro"` | no |
| zone | Zone to host vm in | `string` | `"us-central1-a"` | no |
| bucket\_name | Name of google bucket | `string` | `"mflinn-infratest-bucket"` | no |
| bucket\_location | location to host the bucket | `string` | `"US"` | no |

## Outputs

| Name | Description |
|------|-------------|
| instance\_name | n/a |
| public\_up | n/a |
| bucket\_url | n/a |

