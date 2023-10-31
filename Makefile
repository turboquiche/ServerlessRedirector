build:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 GOFLAGS=-trimpath go build -tags lambda.norpc -mod=readonly -ldflags='-s -w' -o bootstrap redirector/main.go

clean:
	rm -rf bootstrap .serverless/

deploy: clean build
	sls deploy --verbose

remove:
	sls remove --verbose