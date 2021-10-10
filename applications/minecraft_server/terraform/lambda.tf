module "cmd_instance_start" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "cmd_instance_start"
  description   = "starts the minecraft server ec2 instance"
  handler       = "main.go"
  runtime       = "go1.x"
  publish       = true

  source_path = "${local.lambda.root}/cmd/instance/start"

  store_on_s3 = true
  s3_bucket   = aws_s3_bucket.lambda.id

  environment_variables = {
    AWS_REGION                   = "us-west-2"
    MINECRAFT_SERVER_INSTANCE_ID = module.minecraft_server.id
  }
}

module "cmd_instance_stop" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "cmd_instance_stop"
  description   = "stops the minecraft server ec2 instance"
  handler       = "main.go"
  runtime       = "go1.x"
  publish       = true

  source_path = "${local.lambda.root}/cmd/instance/stop"

  store_on_s3 = true
  s3_bucket   = aws_s3_bucket.lambda.id

  environment_variables = {
    AWS_REGION                   = "us-west-2"
    MINECRAFT_SERVER_INSTANCE_ID = module.minecraft_server.id
  }
}

module "cmd_instance_uptime" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "cmd_instance_uptime"
  description   = "gets uptime details for the minecraft server ec2 instance"
  handler       = "main.go"
  runtime       = "go1.x"
  publish       = true

  source_path = "${local.lambda.root}/cmd/instance/uptime"

  store_on_s3 = true
  s3_bucket   = aws_s3_bucket.lambda.id

  environment_variables = {
    AWS_REGION                   = "us-west-2"
    MINECRAFT_SERVER_INSTANCE_ID = module.minecraft_server.id
  }
}

module "cmd_instance_ip" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "cmd_instance_ip"
  description   = "gets the ip of the minecraft server ec2 instance"
  handler       = "main.go"
  runtime       = "go1.x"
  publish       = true

  source_path = "${local.lambda.root}/cmd/instance/ip"

  store_on_s3 = true
  s3_bucket   = aws_s3_bucket.lambda.id

  environment_variables = {
    AWS_REGION                   = "us-west-2"
    MINECRAFT_SERVER_INSTANCE_ID = module.minecraft_server.id
  }
}
