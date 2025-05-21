module github.com/m11ano/mipt-webdev-course/backend/services/orders

go 1.23.3

require github.com/m11ano/mipt-webdev-course/backend/protos v0.0.0

replace github.com/m11ano/mipt-webdev-course/backend/protos => ../../protos

require github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc v0.0.0

replace github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc => ../../clients/clgrpc

require github.com/m11ano/mipt-webdev-course/backend/temporal-app v0.0.0

replace github.com/m11ano/mipt-webdev-course/backend/temporal-app => ../../temporal-app

replace github.com/m11ano/mipt-webdev-course/backend/services/auth/pkg/auth => ../auth/pkg/auth

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2 v2.0.0
	github.com/avito-tech/go-transaction-manager/trm/v2 v2.0.0
	github.com/georgysavva/scany/v2 v2.1.4
	github.com/go-playground/validator/v10 v10.26.0
	github.com/gofiber/fiber/v2 v2.52.6
	github.com/gofiber/swagger v1.1.1
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.2
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/imperatorofdwelling/Website-backend v0.0.0-20240718064027-77c56fad23ad
	github.com/jackc/pgx/v5 v5.7.5
	github.com/m11ano/e v1.0.3
	github.com/m11ano/mipt-webdev-course/backend/clients/auth v0.0.0-20250521131324-52f118631ed0
	github.com/pressly/goose/v3 v3.24.3
	github.com/samber/lo v1.50.0
	github.com/shopspring/decimal v1.4.0
	github.com/stretchr/testify v1.10.0
	github.com/swaggo/swag v1.16.4
	go.temporal.io/sdk v1.34.0
	go.uber.org/fx v1.24.0
	golang.org/x/sync v0.14.0
	google.golang.org/grpc v1.72.1
)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/nexus-rpc/sdk-go v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/swaggo/files/v2 v2.0.2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.62.0 // indirect
	go.temporal.io/api v1.46.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250512202823-5a2f75b736a9 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250512202823-5a2f75b736a9 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
