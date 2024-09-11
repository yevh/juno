data "aws_iam_policy_document" "ssm_lifecycle_trust" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["events.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ssm_lifecycle" {
  name               = "SSMLifecycle"
  assume_role_policy = data.aws_iam_policy_document.ssm_lifecycle_trust.json
}

resource "aws_iam_policy" "ssm_lifecycle" {
  name   = "SSMLifecycle"
  policy = jsonencode(
    {
        Version = "2012-10-17"
        Statement = [
            {
            Action = "ssm:SendCommand"
            Condition = {
                StringEquals = {
                "ec2:ResourceTag/Project" = "juno"
                }
            }
            Effect   = "Allow"
            Resource = "*"
            },
        ],
    }
  )
}

resource "aws_iam_role_policy_attachment" "ssm_lifecycle" {
  policy_arn = aws_iam_policy.ssm_lifecycle.arn
  role       = aws_iam_role.ssm_lifecycle.name
}

resource "aws_ssm_document" "stop_instance" {
  name          = "stop_instance"
  document_type = "Command"

  content = jsonencode({
    schemaVersion = "1.2"
    description   = "Stop an instance"
    parameters    = {}
    runtimeConfig = {
      "aws:runShellScript" = {
        properties = [
          {
            id         = "0.aws:runShellScript"
            runCommand = ["shutdown -H now"]
          }
        ]
      }
    }
  })
}

resource "aws_cloudwatch_event_rule" "stop_instances" {
  name                = "StopInstance"
  description         = "Stop instances nightly"
  schedule_expression = "cron(0 3 * * ? *)"
}

resource "aws_cloudwatch_event_target" "stop_instances" {
  target_id = aws_cloudwatch_event_rule.stop_instances.name
  arn       = aws_ssm_document.stop_instance.arn
  rule      = aws_cloudwatch_event_rule.stop_instances.name
  role_arn  = aws_iam_role.ssm_lifecycle.arn

  run_command_targets {
    key    = "tag:Project"
    values = ["juno"]
  }
}