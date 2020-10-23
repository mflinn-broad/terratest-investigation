output "instance_name" {
  value = google_compute_instance.test_vm.name
}

output "public_ip" {
  value = google_compute_instance.test_vm.network_interface[0].access_config[0].nat_ip
}

output "bucket_url" {
  value = google_storage_bucket.test_bucket.url
}
