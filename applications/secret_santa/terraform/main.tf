terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.62.0"
    }
  }
  backend "s3" {
    bucket = "com.reubenreyes.tfstate"
    key    = "applications/secret_santa"
    region = "us-west-2"
  }
}

provider "aws" {
  region = "us-west-2"
  default_tags {
    tags = {
      Account     = "hub"
      Application = "secret_santa"
      Prefix      = "ss2021"
      ManagedBy   = "applications/secret_santa"
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

variable "discord_channel_id" {
  description = "Discord channel on which to notify"
  type        = string
  default = "735641786306920471"
}

variable "draw_url" {
  description = "URL of the DrawNames page"
  type        = string
}

locals {
  metadata = {
    prefix = "ss2021"
  }
  git = {
    root = "${path.root}/../../../../"
  }
}
