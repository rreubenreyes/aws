resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_internet_gateway" "main_public_gateway" {
  vpc_id = aws_vpc.main.id
}
