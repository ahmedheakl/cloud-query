init:
	cd terraform && terraform init

apply:
	cd api/ && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main . 
	cd terraform/ && terraform apply -auto-approve

clean:
	rm api/main.zip api/main

# terraform destroy requires the zip file to exist
destroy:
	cd api/ && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main . 
	cd terraform/ && terraform destroy -auto-approve
	make clean
