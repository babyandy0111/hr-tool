windows:
	 CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go

mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go