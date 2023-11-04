// SMA
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.26.0 -g application -p indicators -i ./../../../api/asyncapi/indicators.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.26.0 -g user        -p indicators -i ./../../../api/asyncapi/indicators.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.26.0 -g types       -p indicators -i ./../../../api/asyncapi/indicators.yaml -o ./types.gen.go

package indicators
