Run the service locally: `go run cmd/pets/main.go --config=./config/app/local.yaml`

Regenerate proto code (from root): 
```
protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
proto/v1beta1/pets/pets_api.proto
```