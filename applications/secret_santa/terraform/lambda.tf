locals {
  lambda = {
    root   = "${path.root}/../lambda"
    prefix = "mc"
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

  function_name = "${local.lambda.prefix}_start_instance"
  description   = "reminds people to sign up for secret santa"
  handler       = "main"
  runtime       = "go1.x"
  publish       = true

  source_path = [{
    path     = "${local.lambda.root}/remind"
    commands = local.lambda.deployment.source_path.commands
  }]

  store_on_s3 = true
  s3_bucket   = local.lambda.deployment.s3.bucket
  s3_prefix   = local.lambda.deployment.s3.prefix

  attach_policy_jsons    = true
  number_of_policy_jsons = 2
  policy_jsons = [
    data.aws_iam_policy_document.describe_instances.json,
    data.aws_iam_policy_document.start_minecraft_server_instance.json,
  ]

  environment_variables = {
    DISCORD_TOKEN_SECRET_ID   = aws_secretsmanager_secret.discord_token.id
    DYNAMO_TABLE_PARTICIPANTS = aws_dynamodb_table.participants.id
  }
}