resource "aws_secretsmanager_secret" "discord_token" {
  name = "applications/secret_santa/discord_token"
}
