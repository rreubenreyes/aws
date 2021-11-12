data "aws_iam_policy_document" "read_write_participant" {
  statement {
    sid       = "ReadWriteParticipant"
    actions   = ["dynamodb:GetItem", "dynamodb:UpdateItem"]
    resources = [aws_dynamodb_table.participants.arn]
  }
}

data "aws_iam_policy_document" "list_static_bucket" {
  statement {
    sid       = "ListStaticBucket"
    actions   = ["s3:ListBucket"]
    resources = [aws_s3_bucket.static.arn]
  }
}

data "aws_iam_policy_document" "read_participant_photos" {
  statement {
    sid       = "ReadParticipantPhotos"
    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.static.arn}/participant_photos/*"]
  }
}

data "aws_iam_policy_document" "read_discord_token_secret" {
  statement {
    sid       = "ReadDiscordTokenSecret"
    actions   = ["secretsmanager:GetSecretValue"]
    resources = [aws_secretsmanager_secret.discord_token.arn]
  }
}
