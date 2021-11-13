resource "aws_cloudwatch_event_rule" "daily_remind" {
  name        = "${local.metadata.prefix}_daily_remind"
  description = "remind people to sign up for secret santa once a day"
  is_enabled  = true

  schedule_expression = "rate(1 day)"
}

resource "aws_cloudwatch_event_target" "daily_remind" {
  arn  = module.remind.lambda_function_arn
  rule = aws_cloudwatch_event_rule.daily_remind.id
}
