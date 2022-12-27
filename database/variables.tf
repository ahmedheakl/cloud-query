variable "username" {
  type        = string
  default     = "postgres"
  sensitive   = true
  description = "username used when logging into the database"
}

variable "password" {
  type        = string
  default     = "postgres"
  sensitive   = true
  description = "password used when logging into the database"
}
