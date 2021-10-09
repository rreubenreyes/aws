# users
## users
resource "aws_iam_user" "reuben" {
  name = "reuben"
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
