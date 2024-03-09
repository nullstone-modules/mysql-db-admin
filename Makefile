NAME := mysql-db-admin

.PHONY: tools build

tools:
	go install github.com/aws/aws-lambda-go/cmd/build-lambda-zip@latest

build:
	mkdir -p ./aws/tf/files
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o -tags lambda.norpc ./aws/tf/files/bootstrap ./aws/
	# Run build on gcp to ensure a successful build, we discard it
	GOOS=linux GOARCH=amd64 go build -o ./gcp/tf/files/mysql-db-admin ./gcp/; rm -f ./gcp/tf/files/mysql-db-admin

package: tools
	# Package aws module using build-lambda-zip which produces a viable package from any OS
	cd ./aws/tf && build-lambda-zip --output files/mysql-db-admin.zip files/bootstrap
	# Package gcp module (source code instead of binary)
	# For GCP, main.go *must* be in the root of the zip file
	cp gcp/main.go main.go && \
		zip -r gcp/tf/files/mysql-db-admin.zip go.mod go.sum main.go ./api/ ./mysql/ ./vendor/; \
		rm main.go

acc: acc-up acc-run acc-down

acc-up:
	cd acc && docker-compose -p mysql-db-admin-acc up -d db

acc-run:
	ACC=1 gotestsum ./acc/...

acc-down:
	cd acc && docker-compose -p mysql-db-admin-acc down
