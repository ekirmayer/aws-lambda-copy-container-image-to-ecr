resource "aws_ecr_repository" "ecr" {
  name                 = "lambda/go/tester"
  image_tag_mutability = "MUTABLE"

}