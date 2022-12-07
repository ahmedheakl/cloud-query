output "crawler-name" {
  value       = aws_glue_crawler.glue-crawler.name
  description = "Name of the created crawler"
}

output "bucket-name" {
  value       = aws_s3_bucket.s3_bucket.id
  description = "Name of the created S3 bucket"
}

output "athena-catalog" {
  value       = "AwsDataCatalog"
  description = "Athena data catalog where the database is stored"
}

output "athena-database" {
  value       = aws_glue_catalog_database.glue-db.name
  description = "Athena database containing the table"
}

output "athena-table" {
  value       = var.s3_target_folder
  description = "Athena table to be queried"
}

output "api-url" {
  value       = module.queries-api.lambda_function_url
  description = "API URL"
}