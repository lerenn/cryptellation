// SMA
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.30.2 -g application -p asyncapi -i ../asyncapi.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.30.2 -g user        -p asyncapi -i ../asyncapi.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.30.2 -g types       -p asyncapi -i ../asyncapi.yaml -o ./types.gen.go

package asyncapi
