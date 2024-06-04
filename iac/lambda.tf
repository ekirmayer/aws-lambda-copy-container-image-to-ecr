module "lambda" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.4.0"

  function_name        = "copy-image"
  description          = "Copy image to ECR based of SQS messages"
  handler              = "main.handler"
  timeout              = 120
  create_package       = false
  package_type         = "Image"
  architectures        = ["x86_64"]
  image_uri            = "${data.aws_caller_identity.this.account_id}.dkr.ecr.${data.aws_region.this.name}.amazonaws.com/lambda/go/tester:test.4"
  image_config_command = ["main.handler"]

  # role
  role_name   = "copy-image-to-ecr"
  policy_path = "/test/"
  role_path   = "/test/"

  environment_variables = {
    "TRIGGER": "SQS"
  }
  event_source_mapping = {
    sqs = {
      event_source_arn        = aws_sqs_queue.sqs.arn
      function_response_types = ["ReportBatchItemFailures"]
      batch_size = 1
      scaling_config = {
        maximum_concurrency = 20
      }
    }
  }
  create_current_version_allowed_triggers = false

  allowed_triggers = {
    sqs = {
      principal  = "sqs.amazonaws.com"
      source_arn = aws_sqs_queue.sqs.arn
    }
  }
  attach_policy_statements = true
  policy_statements = {
    ecr = {
      effect    = "Allow",
      actions   = [
        "ecr:BatchGetImage",
        "ecr:BatchCheckLayerAvailability",
				"ecr:BatchDeleteImage",
        "ecr:CompleteLayerUpload",
        "ecr:DescribeImages",
        "ecr:GetAuthorizationToken",
        "ecr:BatchGetImage",
        "ecr:GetDownloadUrlForLayer",
				"ecr:UploadLayerPart",
				"ecr:InitiateLayerUpload",
        "ecr:ListImages",
				"ecr:PutImage"],
      resources = ["*"]
    }
  }

  attach_policies    = true
  number_of_policies = 1

  policies = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"
  ]
}
