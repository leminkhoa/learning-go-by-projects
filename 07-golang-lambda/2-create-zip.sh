GOOS=linux GOARCH=arm64 go build -o bootstrap -tags lambda.norpc main.go
zip myFunction.zip bootstrap
