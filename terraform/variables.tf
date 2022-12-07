variable "s3_target_bucket" {
  type        = string
  default     = "layoffs-data"
  description = "Name of the S3 bucket where the data will be available"
  nullable    = false
}

variable "s3_target_folder" {
  type        = string
  default     = "data"
  description = "Name of the directory where the csv will be stored, this is also used to be the name of the table"
  nullable    = false
}

variable "s3_target_file" {
  type        = string
  description = "Name of file to be uploaded for querying"
  default     = "layoffs_data.csv"
  nullable    = false
}

variable "target_file_path" {
  type        = string
  description = "Path of file to be uploaded for querying"
  default     = "../layoffs_data.csv"
  nullable    = false
}

variable "athena_db" {
  type        = string
  default     = "layoffsdb"
  description = "Name of the database created for querying"
  validation {
    condition     = lower(var.athena_db) == var.athena_db
    error_message = "Name must be all lowercase"
  }
  nullable = false
}

variable "crawler-name" {
  type        = string
  default     = "s3-crawler"
  description = "Name of the Glue Crawler"
}