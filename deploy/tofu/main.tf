resource "google_cloud_run_v2_service" "kycelisd" {
  project  = var.GOOGLE_PROJECT_ID
  name     = "kycelisd"
  location = "us-central1"
  deletion_protection = false
  ingress = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = var.CONTAINER_IMAGE
      resources {
        limits = {
          cpu    = "10m"
          memory = "128Mi"
        }
      }
    }
  }
}
