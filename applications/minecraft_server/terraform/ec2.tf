locals {
  ec2 = {
    ami = {
      ubuntu_2004_lts_x86 = "ami-03d5c68bab01f3496"
    }
  }
}

module "minecraft_server" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "3.2.0"

  name                    = "minecraft_server"
  ami                     = local.ec2.ami.ubuntu_2004_lts_x86
  instance_type           = "m5a.large"
  disable_api_termination = true
}

