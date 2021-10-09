module "minecraft_server" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "3.2.0"

  name          = "minecraft_server"
  instance_type = "m5a.large"
}
