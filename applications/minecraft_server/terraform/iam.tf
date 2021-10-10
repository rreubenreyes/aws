data "aws_iam_policy_document" "describe_instances" {
  statement {
    sid       = "DescribeInstances"
    actions   = ["ec2:DescribeInstances"]
    resources = ["*"]
  }
}

data "aws_iam_policy_document" "start_minecraft_server_instance" {
  statement {
    sid       = "StartMinecraftServerInstance"
    actions   = ["ec2:StartInstances"]
    resources = [module.minecraft_server.arn]
  }
}

data "aws_iam_policy_document" "stop_minecraft_server_instance" {
  statement {
    sid       = "StopMinecraftServerInstance"
    actions   = ["ec2:StopInstances"]
    resources = [module.minecraft_server.arn]
  }
}
