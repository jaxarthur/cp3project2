build :	main.go
		env GOOS=windows GOARCH=amd64 go build -o project2win.exe
		env GOOS=darwin GOARCH=amd64 go build -o project2mac
		env GOOS=linux GOARCH=amd64 go build -o project2linux