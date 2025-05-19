module github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc

go 1.23.3

require github.com/m11ano/mipt-webdev-course/backend/protos v0.0.0

replace github.com/m11ano/mipt-webdev-course/backend/protos => ../../protos

require (
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.2
	github.com/m11ano/avito-pvz v0.0.0-20250414202626-53e162156757
	github.com/samber/lo v1.50.0
	github.com/shopspring/decimal v1.4.0
	google.golang.org/genproto v0.0.0-20250512202823-5a2f75b736a9
	google.golang.org/grpc v1.72.1
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.2 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250505200425-f936aa4a68b2 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
