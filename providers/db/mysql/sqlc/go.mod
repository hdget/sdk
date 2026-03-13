module github.com/hdget/sdk/providers/db/mysql/sqlc

go 1.24.0

replace github.com/hdget/sdk/common => ../../../../common

require (
	github.com/go-sql-driver/mysql v1.9.3
	github.com/hdget/sdk/common v0.0.7
	github.com/pkg/errors v0.9.1
	go.uber.org/fx v1.24.0
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	golang.org/x/exp v0.0.0-20251219203646-944ab1f22d93 // indirect
	golang.org/x/sys v0.39.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251213004720-97cd9d5aeac2 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
