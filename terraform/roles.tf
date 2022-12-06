resource "aws_iam_policy" "glue-policy" {
  name = "glue-policy"
  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Action" : [
          "glue:*",
          "athena:*",
          "s3:*",
          "ec2:*",
          "iam:*",
          "cloudwatch:*",
          "logs:*"
        ],
        "Resource" : [
          "*"
        ]
      }
    ]
  })
}

resource "aws_iam_role" "glue-role" {
  name = "glue-role"
  assume_role_policy = jsonencode(
    {
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Sid" : "",
          "Effect" : "Allow",
          "Principal" : {
            "Service" : [
              "glue.amazonaws.com"
            ]
          },
          "Action" : "sts:AssumeRole"
        }
      ]
    }
  )
}

resource "aws_iam_role_policy_attachment" "glue-policy-attach" {
  role       = aws_iam_role.glue-role.name
  policy_arn = aws_iam_policy.glue-policy.arn
}

