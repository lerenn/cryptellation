package ci

import "dagger.io/dagger"

func otelCollector(client *dagger.Client) *dagger.Container {
	config := client.Host().File("./tools/config/ci/otel-collector.yaml")

	return client.Container().
		From("otel/opentelemetry-collector-contrib:0.88.0").
		WithFile("/etc/otelcol-contrib/config.yaml", config).
		WithoutUser().
		WithExposedPort(4317).
		WithExposedPort(4318)
}

func Uptrace(client *dagger.Client) (uptrace, otelcollector *dagger.Service) {
	config := client.Host().File("./tools/config/ci/uptrace.yaml")

	uptrace = client.Container().
		From("uptrace/uptrace:1.6.2").
		WithFile("/etc/uptrace/uptrace.yml", config).
		WithServiceBinding("postgres", uptracePostgres(client)).
		WithServiceBinding("clickhouse", uptraceClickHouse(client)).
		WithExposedPort(4317).
		WithExposedPort(4318).
		AsService()

	otelcollector = otelCollector(client).
		WithServiceBinding("uptrace", uptrace).
		AsService()

	return
}

func uptracePostgres(client *dagger.Client) *dagger.Service {
	return client.Container().
		From("postgres:15-alpine").
		WithEnvVariable("POSTGRES_USER", "uptrace").
		WithEnvVariable("POSTGRES_PASSWORD", "uptrace").
		WithEnvVariable("POSTGRES_DB", "uptrace").
		WithExposedPort(5432).
		AsService()
}

func uptraceClickHouse(client *dagger.Client) *dagger.Service {
	return client.Container().
		From("clickhouse/clickhouse-server:23.7").
		WithEnvVariable("CLICKHOUSE_DB", "uptrace").
		WithExposedPort(8123).
		WithExposedPort(9000).
		AsService()
}
