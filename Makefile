test:
	@go test ./...

deps:
	@go mod download

build: deps
	@go build -o metricssvc cmd/metrics/main.go

run-without-swagger: build
	./metricssvc -withswagger=false

run: build
	./metricssvc -withswagger=true -swaggerdir=${SWAGGER_UI}