locals {
  s3 = {
    bucket = "${data.terraform_remote_state.hub.outputs.aws.s3.namespace_prefix}.applications.minecraft-server"
  }
}

resource "aws_s3_bucket" "deployment_artifacts" {
  bucket = local.s3.bucket
  acl    = "private"
}

resource "aws_s3_bucket_public_access_block" "deployment_artifacts" {
  bucket = aws_s3_bucket.deployment_artifacts.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
