resource "aws_dynamodb_table" "participants" {
  name     = "${local.metadata.prefix}_participants"
  hash_key = "name"

  attribute {
    name = "name"
    type = "S"
  }
}
