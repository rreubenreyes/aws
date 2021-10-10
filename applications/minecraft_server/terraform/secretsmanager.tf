resource "aws_secretsmanager_secret" "ec2_instance_private_key" {
  name = "applications/minecraft_server/ec2_instance/private_key"
}
