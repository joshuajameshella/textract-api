lambda:
	rm -rf build/
	GOOS=linux go build -o build/main cmd/main.go
	zip -jrm build/main.zip build/main

