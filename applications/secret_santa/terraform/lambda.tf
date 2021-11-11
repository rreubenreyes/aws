locals {
  lambda = {
    root = "${path.root}/../lambda"
    deployment = {
      s3 = {
        bucket = aws_s3_bucket.deployment_artifacts.id
        prefix = "lambda"
      }
      source_path = {
        commands = [
          "GOARCH=amd64 GOOS=linux go build main.go",
          ":zip main"
        ]
      }
    }
  }
}

module "remind" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "${local.metadata.prefix}_remind"
  description   = "reminds people to sign up for secret santa"
  handler       = "main"
  runtime       = "go1.x"
  publish       = true
  timeout       = 60

  source_path = [{
    path     = "${local.lambda.root}/remind"
    commands = local.lambda.deployment.source_path.commands
  }]

  store_on_s3 = true
  s3_bucket   = local.lambda.deployment.s3.bucket
  s3_prefix   = local.lambda.deployment.s3.prefix

  attach_policy_jsons    = true
  number_of_policy_jsons = 1
  policy_jsons = [
    data.aws_iam_policy_document.get_participant.json,
  ]

  environment_variables = {
    DISCORD_CHANNEL_ID          = var.discord_channel_id
    DISCORD_TOKEN_SECRET_ID     = aws_secretsmanager_secret.discord_token.id
    DRAW_URL                    = var.draw_url
    DYNAMODB_TABLE_PARTICIPANTS = aws_dynamodb_table.participants.id
  }
}
