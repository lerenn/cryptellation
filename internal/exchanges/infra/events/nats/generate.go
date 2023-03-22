//go:generate asyncapi-codegen -g application -p generated -i ../../../../../api/asyncapi-spec/exchanges.yaml -o ./generated/app.gen.go
//go:generate asyncapi-codegen -g client      -p generated -i ../../../../../api/asyncapi-spec/exchanges.yaml -o ./generated/client.gen.go
//go:generate asyncapi-codegen -g broker      -p generated -i ../../../../../api/asyncapi-spec/exchanges.yaml -o ./generated/broker.gen.go
//go:generate asyncapi-codegen -g types       -p generated -i ../../../../../api/asyncapi-spec/exchanges.yaml -o ./generated/types.gen.go
//go:generate asyncapi-codegen -g nats        -p generated -i ../../../../../api/asyncapi-spec/exchanges.yaml -o ./generated/nats.gen.go

package generated
