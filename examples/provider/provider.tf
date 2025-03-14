# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    devlake = {
      source = "registry.terraform.io/incubator-devlake-terraform/devlake"
    }
  }
}

provider "devlake" {
  host  = "http://localhost:8080"
  token = "whatever"
}
