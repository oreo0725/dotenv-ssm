bin/ssm2dotenv: *.go
	CGO_ENABLED=0 go build -o $@ .
