data "aws_iam_policy_document" "get_participant" {
  statement {
    sid       = "GetParticipant"
    actions   = ["dynamodb:GetItem"]
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

data "aws_iam_policy_document" "write_participant_photos" {
  statement {
    sid       = "WriteParticipantPhotos"
    actions   = ["s3:PutObject"]
    resources = ["${aws_s3_bucket.static.arn}/participant_photos/*"]
  }
}
