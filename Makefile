build:
	@CGO_ENABLED=0 go build -o dist/pkg-diff ./app/services/pkg-diff/

run: build
	@dist/pkg-diff
