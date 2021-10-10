resource "aws_security_group" "minecraft_server" {
  name = "minecraft_server_instance"
}

resource "aws_security_group_rule" "default_egress" {
  security_group_id = aws_security_group.minecraft_server.id

  type      = "egress"
  protocol  = "-1"
  from_port = 0
  to_port   = 0
  cidr_blocks = [
    "0.0.0.0/0"
  ]
}

resource "aws_security_group_rule" "default_ingress" {
  security_group_id = aws_security_group.minecraft_server.id

  type      = "ingress"
  protocol  = "-1"
  from_port = 0
  to_port   = 0
  self      = true
}

resource "aws_security_group_rule" "allow_minecraft_server_ingress" {
  security_group_id = aws_security_group.minecraft_server.id

  type      = "ingress"
  protocol  = "tcp"
  from_port = 25565
  to_port   = 25565
  cidr_blocks = [
    "0.0.0.0/0"
  ]
}

resource "aws_security_group_rule" "allow_ssh_ingress" {
  security_group_id = aws_security_group.minecraft_server.id

  type = "ingress"

  protocol  = "tcp"
  from_port = 22
  to_port   = 22
  cidr_blocks = [
    "0.0.0.0/0"
  ]
}
