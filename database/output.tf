output "address" {
  value       = aws_db_instance.default.address
  description = "Address of the db"
  sensitive = true
}

output "name" {
  value = aws_db_instance.default.db_name
  sensitive = true
}

output "username" {
  value     = var.username
  sensitive = true
}

output "password" {
  value     = var.password
  sensitive = true
}
