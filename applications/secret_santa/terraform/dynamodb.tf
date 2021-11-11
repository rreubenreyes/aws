resource "aws_dynamodb_table" "participants" {
  name     = "${local.metadata.prefix}_participants"
  hash_key = "name"

  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "name"
    type = "S"
  }
}
