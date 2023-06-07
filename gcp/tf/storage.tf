locals {
  location         = startswith(local.region, "us") ? "us" : (startswith(local.region, "eu") ? "eu" : "asia")
  package_filename = "${path.module}/files/mysql-db-admin.zip"
}

resource "google_storage_bucket" "binaries" {
  name          = "${var.name}-binaries"
  location      = local.location
  labels        = var.labels
  force_destroy = true
  storage_class = "MULTI_REGIONAL"
}

resource "google_storage_bucket_object" "binary" {
  bucket = google_storage_bucket.binaries.name
  name   = "mysql-db-admin.zip"
  source = local.package_filename
}
