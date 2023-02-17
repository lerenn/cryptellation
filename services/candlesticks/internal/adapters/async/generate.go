//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g application -p async -i ../../../../../api/asyncapi-spec/candlesticks.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g client      -p async -i ../../../../../api/asyncapi-spec/candlesticks.yaml -o ./client.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g broker      -p async -i ../../../../../api/asyncapi-spec/candlesticks.yaml -o ./broker.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g types       -p async -i ../../../../../api/asyncapi-spec/candlesticks.yaml -o ./types.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g nats        -p async -i ../../../../../api/asyncapi-spec/candlesticks.yaml -o ./nats.gen.go

package async
