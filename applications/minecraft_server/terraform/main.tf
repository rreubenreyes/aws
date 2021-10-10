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

data "terraform_remote_state" "hub" {
  backend = "s3"
  config = {
    bucket = "com.reubenreyes.tfstate"
    key    = "hub"
    region = "us-west-2"
  }
}

locals {
  git = {
    root = "${path.root}/../../../../"
  }
  ec2 = {
    ami = {
      ubuntu_2004_lts_x86 = "ami-03d5c68bab01f3496"
    }
  }
  lambda = {
    root = "${path.root}/../lambda"
    prefix = "mc"
    deployment = {
      s3 = {
        bucket = aws_s3_bucket.deployment_artifacts.id
        prefix = "lambda"
      }
    }
  }
  s3 = {
    namespace_prefix = "${data.terraform_remote_state.hub.outputs.aws.s3.namespace_prefix}.applications.minecraft-server"
  }
}
