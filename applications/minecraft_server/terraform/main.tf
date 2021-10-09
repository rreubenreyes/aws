terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.62.0"
    }
  }
  backend "s3" {
    bucket = "com.reubenreyes.tfstate"
    key    = "applications/minecraft_server"
    region = "us-west-2"
  }
}

provider "aws" {
  region = "us-west-2"
  default_tags {
    tags = {
      Account     = "hub"
      Application = "minecraft_server"
      ManagedBy   = "applications/minecraft_server"
    }
  }
}

locals {
  git = {
    root = "${path.root}/../../../../"
  }
}
