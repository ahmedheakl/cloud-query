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
  name        = var.athena_db
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
  provisioner "local-exec" {
    command = "aws glue start-crawler --name ${self.name}"
  }
}