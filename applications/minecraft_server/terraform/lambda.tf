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

module "start_instance" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "${local.lambda.prefix}_start_instance"
  description   = "starts the minecraft server ec2 instance"
  handler       = "main"
  runtime       = "go1.x"
  publish       = true

  source_path = [{
    path     = "${local.lambda.root}/cmd/instance/start"
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
    MINECRAFT_SERVER_INSTANCE_ID = module.minecraft_server.id
  }
}

module "stop_instance" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "${local.lambda.prefix}_stop_instance"
  description   = "stops the minecraft server ec2 instance"
  handler       = "main"
  runtime       = "go1.x"
  publish       = true

  source_path = [{
    path     = "${local.lambda.root}/cmd/instance/stop"
    commands = local.lambda.deployment.source_path.commands
  }]

  store_on_s3 = true
  s3_bucket   = local.lambda.deployment.s3.bucket
  s3_prefix   = local.lambda.deployment.s3.prefix

  attach_policy_jsons    = true
  number_of_policy_jsons = 2
  policy_jsons = [
    data.aws_iam_policy_document.describe_instances.json,
    data.aws_iam_policy_document.stop_minecraft_server_instance.json,
  ]

  environment_variables = {
    MINECRAFT_SERVER_INSTANCE_ID = module.minecraft_server.id
  }
}

module "get_instance_uptime" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "${local.lambda.prefix}_get_instance_uptime"
  description   = "gets uptime details for the minecraft server ec2 instance"
  handler       = "main"
  runtime       = "go1.x"
  publish       = true

  source_path = [{
    path     = "${local.lambda.root}/cmd/instance/uptime"
    commands = local.lambda.deployment.source_path.commands
  }]

  store_on_s3 = true
  s3_bucket   = local.lambda.deployment.s3.bucket
  s3_prefix   = local.lambda.deployment.s3.prefix

  attach_policy_jsons    = true
  number_of_policy_jsons = 1
  policy_jsons = [
    data.aws_iam_policy_document.describe_instances.json,
  ]

  environment_variables = {
    MINECRAFT_SERVER_INSTANCE_ID = module.minecraft_server.id
  }
}

module "get_instance_ip" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "${local.lambda.prefix}_get_instance_ip"
  description   = "gets the ip of the minecraft server ec2 instance"
  handler       = "main"
  runtime       = "go1.x"
  publish       = true

  source_path = [{
    path     = "${local.lambda.root}/cmd/instance/ip"
    commands = local.lambda.deployment.source_path.commands
  }]

  store_on_s3 = true
  s3_bucket   = local.lambda.deployment.s3.bucket
  s3_prefix   = local.lambda.deployment.s3.prefix

  attach_policy_jsons    = true
  number_of_policy_jsons = 1
  policy_jsons = [
    data.aws_iam_policy_document.describe_instances.json,
  ]

  environment_variables = {
    MINECRAFT_SERVER_INSTANCE_ID = module.minecraft_server.id
  }
}
