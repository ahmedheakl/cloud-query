output "crawler-name" {
    value = aws_glue_crawler.glue-crawler.name
}

output "bucket-name" {
    value = aws_s3_bucket.s3_bucket.id
}

output "athena-catalog" {
    value = "AwsDataCatalog"
}

output "athena-database" {
    value = aws_glue_catalog_database.glue-db.name
}

output "athena-table" {
    value = var.s3_target_folder
}