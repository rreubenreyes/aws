terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.62.0"
    }
  }
  backend "s3" {
    bucket = "com.reubenreyes.tfstate"
    key    = "hub"
    region = "us-west-2"
  }
}

provider "aws" {
  region = "us-west-2"
  default_tags {
    tags = {
      Account   = "hub"
      ManagedBy = "hub"
    }
  }
}

locals {
  git = {
    root = "${path.root}/../.."
  }
  aws = {
    s3 = {
      namespace_prefix = "com.reubenreyes"
    }
  }
}
