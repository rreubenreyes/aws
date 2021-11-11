data "aws_iam_policy_document" "get_participant" {
  statement {
    sid       = "GetParticipant"
    actions   = ["dynamodb:GetItem"]
    resources = [aws_dynamodb_table.participants.arn]
  }
}
