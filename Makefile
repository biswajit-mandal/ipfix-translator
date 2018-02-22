build:
	cd tools &&  go run build_datasets.go
	go build -o ipfix-translator main.go
