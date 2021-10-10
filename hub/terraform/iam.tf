# users
## users
resource "aws_iam_user" "reuben" {
  name = "reuben"
  path = "/people/"
}

resource "aws_iam_user" "aaron" {
  name = "aaron"
  path = "/people/"
}

resource "aws_iam_user" "clu" {
  name = "clu"
  path = "/people/"
}

resource "aws_iam_user" "ac" {
  name = "ac"
  path = "/people/"
}

resource "aws_iam_user" "steven" {
  name = "steven"
  path = "/people/"
}

resource "aws_iam_user" "ian" {
  name = "ian"
  path = "/people/"
}

## user policy attachments
resource "aws_iam_user_policy_attachment" "admins" {
  for_each = toset([
    aws_iam_user.reuben.name
  ])

  user       = each.value
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}
