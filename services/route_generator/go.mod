module github.com/AlexeyBobkovDev/tripmate/services/route_generator

go 1.27.4

require (
	github.com/AlexeyBobkovDev/tripmate/contracts v0.0.0-20261115224043-2f62ace196f6
	github.com/google/uuid v1.6.0
	github.com/jackc/pgerrcode v0.0.0-20250907135507-afb5586c32a6
	github.com/jackc/pgx/v5 v5.13.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/stretchr/testify v1.12.1
	go.uber.org/mock v0.6.0
	go.uber.org/zap v1.28.0
	google.golang.org/grpc v1.85.0
	google.golang.org/protobuf v1.36.12
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.12-20260825204119-511051f7f437.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.6 // indirect
	golang.org/x/net v0.63.0 // indirect
	golang.org/x/sync v0.25.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
	golang.org/x/text v0.44.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261116203732-fd2c7689a808 // indirect
)

replace github.com/AlexeyBobkovDev/tripmate/contracts => ../../contracts
