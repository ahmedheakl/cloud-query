provider "aws" {
  region = "me-south-1"
}

resource "aws_s3_bucket" "s3_bucket" {
  bucket        = var.s3_target_bucket
  force_destroy = true
}

resource "aws_s3_object" "folder" {
  bucket = aws_s3_bucket.s3_bucket.id
  key    = format("%s/%s", var.s3_target_folder, var.s3_target_file)
  source = format("./%s", var.target_file_path)
}

resource "aws_glue_catalog_database" "glue-db" {
  name = var.athena_db
}

resource "aws_glue_crawler" "glue-crawler" {
  database_name = aws_glue_catalog_database.glue-db.name
  name          = var.crawler-name
  description   = "A crawler to provide data to Athena"
  role          = aws_iam_role.glue-role.arn

  s3_target {
    path = format("s3://%s/%s/", aws_s3_bucket.s3_bucket.id, var.s3_target_folder)
  }

  // start crawling on creation
  // terraform doesn't wait for command to finish
  provisioner "local-exec" {
    command = "aws glue start-crawler --name ${self.name}"
  }
}

data "archive_file" "go_package" {
  type        = "zip"
  source_file = "../api/main"
  output_path = "../api/main.zip"
}

module "queries-api" {
  depends_on                 = [data.archive_file.go_package]
  source                     = "terraform-aws-modules/lambda/aws"
  description                = "API to execute read queries"
  function_name              = "queries-api"
  handler                    = "main"
  runtime                    = "go1.x"
  create_package             = false
  local_existing_package     = "../api/main.zip"
  create_lambda_function_url = true
  attach_policy              = true
  policy                     = aws_iam_policy.api-policy.arn
  timeout                    = 300
  memory_size                = 1024
}