Based on: https://blog.xpnsec.com/aws-lambda-redirector/

## TLDR:
	1. Configure AWS credentials in serverless:
		`sls config credentials --key {access_key} --secret {secret_key} --provider aws -o`
	2. Change teamserver domain/ip in 'serverless.yml'
	3. Build the binary:
		(nix) Binary is built during make deploy
		(win) See long version step 4
	4. Deploy:
		(nix) `make deploy`
		(win) `sls deploy`
	5. Tear Down:
		(nix) `make remove`
		(win) `sls remove`

## Long Version:
	1. Dependencies:
		# Go (nix & win): 
			https://go.dev/doc/install
			(nix) `rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz`
			(nix) `export PATH=$PATH:/usr/local/go/bin`

		# Node (win only): https://nodejs.org/en/download

		# Serverless: 
			(nix) `curl -o- -L https://slss.io/install | VERSION=3.35.2 bash`
			(win) `npm install -g serverless`

	2. Configure AWS credentials in serverless:
		# In AWS:
			1. Create an AWS user:
				> https://console.aws.amazon.com/iamv2/home#/users
			2. Create keys for the AWS user:
				> Click user in table
				> Go to "Security credentials" tab
				> Scroll down to "Access kes"
				> Create access key
				> Copy access key id and secret access key

		# In CLI:
			(nix & win) `sls config credentials`
				> Enter AWS access key id and secret access key
					(nix) `sls config credentials --key {access_key} --secret {secret_key} --provider aws -o`

	3. Set TEAMSERVER parameter in 'serverless.yml':
		Use the domain name or IP associated with your teamserver
		**NOTE** Be sure to include https:// or http:// or else it panics and dies for some reason lol

		Requests made to the lambda endpoint *MUST BE* over HTTPS and port 443
		If using a listener with SSL over a port other than 443, format as https://domain.com:80

		You must also ensure your beacon profile calls out to port 443 to your lambda domain
		The teamserver listener can be any port as long as you set the redirect url as shown above

	4. Build the bootstrap binary:
		(nix) Binary is built during make deploy
		(win) Open cmd in the project folder & run:
			`set GOARCH=amd64`
			`set GOOS=linux`
			`set CGO_ENABLED=0`
			`set GOFLAGS=-trimpath`
			`go build -tags lambda.norpc -mod=readonly -ldflags="-s -w" -o bootstrap redirector/main.go`

	5. Deploy redirector:
		(nix) `make deploy`
		(win) `sls deploy`

	6. Destroy redirector:
		(nix) `make remove`
		(win) `sls remove`
