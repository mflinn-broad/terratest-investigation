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
}

terraform {
  required_version = ">= 0.12.26"
}
