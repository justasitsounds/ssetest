application:
	GOOS=linux GOARCH=amd64 go build -o bin/application application.go

bundle:
	zip bin/application