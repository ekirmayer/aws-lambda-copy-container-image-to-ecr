resource "aws_ecr_repository" "ecr" {
  name                 = "lambda/go/tester"
  image_tag_mutability = "IMMUTABLE"
image_scanning_configuration {
  scan_on_push = false
}
}