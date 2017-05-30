fmt:
	go fmt ./...

lib: fmt
	go build -buildmode=c-shared -o build/libkatana.so ./capi
