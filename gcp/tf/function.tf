resource "google_cloudfunctions2_function" "function" {
  name        = var.name
  location    = local.region
  description = "${var.name} MySQL DB Admin"
  labels      = var.labels

  build_config {
    runtime     = "go120"
    entry_point = "mysql-db-admin"

    source {
      storage_source {
        bucket = google_storage_bucket.binaries.name
        object = google_storage_bucket_object.binary.name
      }
    }
  }

  service_config {
    service_account_email            = google_service_account.executor.email
    available_cpu                    = "4"
    available_memory                 = "4Gi"
    timeout_seconds                  = 20
    max_instance_count               = 100
    max_instance_request_concurrency = 50
    all_traffic_on_latest_revision   = true
    ingress_settings                 = "ALLOW_ALL"
    vpc_connector_egress_settings    = "ALL_TRAFFIC"
    vpc_connector                    = var.vpc_access_connector_name

    secret_environment_variables {
      key        = "DB_CONN_URL"
      project_id = local.project_id
      secret     = google_secret_manager_secret.db_admin_mysql.secret_id
      version    = "latest"
    }
  }
}
