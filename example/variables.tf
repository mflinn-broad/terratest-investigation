variable "google_project" {
  type        = string
  description = "google project to run tests in"
}

variable "instance_name" {
  type        = string
  description = "name of gcp vm instance"
  default     = "terratest-example"
}

variable "machine_type" {
  type        = string
  description = "Machine type of vm"
  default     = "f1-micro"
}

variable "zone" {
  description = "Zone to host vm in"
  type        = string
  default     = "us-central1-a"
}

variable "bucket_name" {
  type        = string
  description = "Name of google bucket"
  default     = "mflinn-infratest-bucket"
}

variable "bucket_location" {
  type        = string
  description = "location to host the bucket"
  default     = "US"
}
