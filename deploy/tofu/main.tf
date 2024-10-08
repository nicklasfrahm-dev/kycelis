resource "google_service_account" "kycelisd" {
  account_id   = "kycelisd"
  display_name = "A service account for the kycelisd service."
}

resource "google_cloud_run_v2_service" "kycelisd" {
  project  = var.GOOGLE_PROJECT_ID
  name     = "kycelisd"
  location = "us-central1"
  deletion_protection = false
  ingress = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.kycelisd.email

    containers {
      image = var.CONTAINER_IMAGE

      liveness_probe {
        failure_threshold     = 3
        initial_delay_seconds = 0
        period_seconds        = 10
        timeout_seconds       = 1

        http_get {
          path = "/health"
          port = 8080
        }
      }
    }
  }
}
