resource "aws_s3_bucket" "tfstate" {
  bucket = "${local.aws.s3.namespace_prefix}.tfstate"
  acl    = "private"
}

resource "aws_s3_bucket_public_access_block" "tfstate" {
  bucket = aws_s3_bucket.tfstate.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
