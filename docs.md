## Terraform

Terraform is an IaaC which uses a language called Hashicorp Configuration Language

Thie means that instead of going into AWS console and setting up the services we need by hand, we write code that describes the services. This helps in managing and collaborating on building the infrastructure. 

To build the infrastructure on your own AWS account you need to:
1. Have an AWS account, with the credentials saved in one of the standard places (`~/.aws/credentials`, as env variables, as variables in .tf file).
2. Download Terraform.
3. `cd` into the terraform folder
4. Modify the environment variables in `secrets.sh`.
5. Run `terraform init`. 
6. Run `source secrets.sh`.

This creates an empty PostgreSQL instance, and sets the database credentials as environemnt variables for the api to read. I'm stil working on loading the data model on DB creation.

The services currently in use are:
1. A single PostgreSQL instance with Public Access.

## API 

The API is written in Golang. You can use the Dockerfile to run the API, or if you want to run it directly:
1. Download and setup golang compiler
2. `cd` into the api folder
3. Run `go mod download` to install dependencies
4. Run `go run .` 

Note that the API connects to the database using the environment variables, so in the case of a restart you will need to run `source secrets.sh` again. Unless the terraform file has changed since the last run, sourcing the script will not change the existing infrastructure.