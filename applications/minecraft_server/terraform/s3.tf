resource "aws_s3_bucket" "deployment_artifacts" {
  bucket = local.s3.namespace_prefix
  acl    = "private"
}

resource "aws_s3_bucket_public_access_block" "deployment_artifacts" {
  bucket = aws_s3_bucket.deployment_artifacts.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
