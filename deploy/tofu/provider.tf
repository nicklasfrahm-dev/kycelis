terraform {
  backend "gcs" {
    bucket = "nicklasfrahm"
    prefix  = "tofu/state/kycelis"
  }
}

provider "google" {
  project = var.GOOGLE_PROJECT_ID

  # This is the most eco-friendly
  # region covered by the free tier.
  region  = "us-central1"
}
