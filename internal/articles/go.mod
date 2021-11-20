module github.com/rezaAmiri123/test-project/internal/articles

go 1.16

replace github.com/rezaAmiri123/test-project/internal/common => ../common

require (
	github.com/go-chi/chi v1.5.4
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/google/uuid v1.3.0
	github.com/gosimple/slug v1.11.2
	github.com/jinzhu/gorm v1.9.16
	github.com/pkg/errors v0.9.1
	github.com/rezaAmiri123/test-project/internal/common v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	google.golang.org/grpc v1.42.0 // indirect
)
