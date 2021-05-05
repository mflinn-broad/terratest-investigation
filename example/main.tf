/**
 * # Terratest Demo Workflow module
 *
 * This is a simple terra form module which create a GCE VM instance and GCS bucket. 
 * It is intended to be used just for a demonstration of a sample automated infrasture 
 * testing workflow using [terratest](https://terratest.gruntwork.io/)
 *
 * To break the automated tests make a pr with the instance or bucket name hard coded
 * rather than using a variable
 *
 * This documentation is generated with [terraform-docs](https://github.com/segmentio/terraform-docs)
 * `terraform-docs markdown --no-sort . > README.md`
 */

// I did it! Now let me Pass!
// "Men are consumed by their desires, untill theres hardly any man left"

resource "google_compute_instance" "test_vm" {
  project      = var.google_project
  name         = var.instance_name
  machine_type = var.machine_type
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1604-lts"
    }
  }

  network_interface {
    network = "default"
    access_config {
    }
  }
}

resource "google_storage_bucket" "test_bucket" {
  project  = var.google_project
  name     = var.bucket_name
  location = var.bucket_location
  uniform_bucket_level_access = true

}

terraform {
  required_version = ">= 0.12.26"
}
