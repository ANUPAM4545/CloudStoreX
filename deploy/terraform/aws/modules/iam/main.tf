resource "aws_iam_role" "app" {
  name = "${var.environment}-cloudstorex-app-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name        = "${var.environment}-cloudstorex-app-role"
    Environment = var.environment
  }
}

resource "aws_iam_policy" "s3_access" {
  name        = "${var.environment}-cloudstorex-s3-policy"
  description = "Policy for CloudStoreX backend to access S3 bucket"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:DeleteObject",
          "s3:ListBucket",
          "s3:GetBucketLocation"
        ]
        Resource = [
          var.s3_bucket_arn,
          "${var.s3_bucket_arn}/*"
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "s3_attach" {
  role       = aws_iam_role.app.name
  policy_arn = aws_iam_policy.s3_access.arn
}
