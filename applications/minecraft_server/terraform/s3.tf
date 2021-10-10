resource "aws_s3_bucket" "lambda" {
  bucket = "${local.s3.namespace_prefix}.lambda"
  acl    = "private"
}

resource "aws_s3_bucket_public_access_block" "lambda" {
  bucket = aws_s3_bucket.lambda.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
