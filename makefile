APP_NAME=go-gin-mongo
MAIN=main.go

# lokasi swag binary (default di GOPATH/bin)
SWAG=$(shell go env GOPATH)/bin/swag

# generate swagger docs
swag:
	@echo "🚀 Generate Swagger Docs..."
	@$(SWAG) init -g $(MAIN)

# run app
run:
	@echo "🚀 Running $(APP_NAME)..."
	go run $(MAIN)

# build binary
build:
	@echo "📦 Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) $(MAIN)
