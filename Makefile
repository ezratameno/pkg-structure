build/pkg-diff:
	@CGO_ENABLED=0 go build -o dist/pkg-diff ./app/services/pkg-diff/

run/pkg-diff: build/pkg-diff
	@dist/pkg-diff


build/pkg-structure:
	@CGO_ENABLED=0 go build -o dist/pkg-structure ./app/services/pkg-structure/

run/pkg-structure: build/pkg-structure
	@dist/pkg-structure