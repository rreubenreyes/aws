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

  source_path = [{
    path     = "${local.lambda.root}/remind"
    commands = local.lambda.deployment.source_path.commands
  }]

  function_name = "${local.metadata.prefix}_remind"
  description   = "reminds people to sign up for secret santa"
  handler       = "main"
  runtime       = "go1.x"
  publish       = true

  memory_size = 512
  timeout     = 60
  environment_variables = {
    DISCORD_CHANNEL_ID          = var.discord_channel_id
    DISCORD_TOKEN_SECRET_ID     = aws_secretsmanager_secret.discord_token.id
    DRAW_URL                    = var.draw_url
    DYNAMODB_TABLE_PARTICIPANTS = aws_dynamodb_table.participants.id
    REGISTER_URL                = var.register_url
    STATIC_BUCKET_ID            = aws_s3_bucket.static.id
  }

  allowed_triggers = {
    eventbridge = {
      principal  = "events.amazonaws.com"
      source_arn = aws_cloudwatch_event_rule.daily_remind.arn
    }
  }

  store_on_s3 = true
  s3_bucket   = local.lambda.deployment.s3.bucket
  s3_prefix   = local.lambda.deployment.s3.prefix

  attach_policy_jsons    = true
  number_of_policy_jsons = 4
  policy_jsons = [
    data.aws_iam_policy_document.read_write_participant.json,
    data.aws_iam_policy_document.list_static_bucket.json,
    data.aws_iam_policy_document.read_participant_photos.json,
    data.aws_iam_policy_document.read_discord_token_secret.json,
  ]
}
