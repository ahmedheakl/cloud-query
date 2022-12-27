provider "aws" {
  region = "me-south-1"
}

# Allow access from any IP over IPv4
resource "aws_security_group" "sg" {
  name        = "allow_global_access"
  description = "Allow global access to RDS"

  ingress {
    to_port     = 5432
    from_port   = 0
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_db_instance" "default" {
  allocated_storage      = 20
  db_name                = "testdb"
  engine                 = "postgres"
  engine_version         = "14.5"
  instance_class         = "db.t3.micro"
  username               = var.username
  password               = var.password
  skip_final_snapshot    = true
  storage_type           = "gp2"
  publicly_accessible    = true
  vpc_security_group_ids = [aws_security_group.sg.id]
}
