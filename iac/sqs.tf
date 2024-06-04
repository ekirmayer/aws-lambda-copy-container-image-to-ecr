resource "aws_sqs_queue" "sqs" {
  name                      = "image-import-queue"
  delay_seconds             = 0
  max_message_size          = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 10
  visibility_timeout_seconds = 130
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.sqs_dlq.arn
    maxReceiveCount     = 4
  })


  tags = {
    Environment = "dev"
  }
}

resource "aws_sqs_queue" "sqs_dlq" {
  name = "image-import-queue-deadletter-queue"
}

resource "aws_sqs_queue_redrive_allow_policy" "sqs_dlq_redrive_allow_policy" {
  queue_url = aws_sqs_queue.sqs_dlq.id

  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue",
    sourceQueueArns   = [aws_sqs_queue.sqs.arn]
  })
}