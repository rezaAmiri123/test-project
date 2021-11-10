module github.com/rezaAmiri123/test-project/internal/users

go 1.16

replace (
	github.com/rezaAmiri123/test-project/internal/common  => ../common
)

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/go-chi/chi v1.5.4
	github.com/go-chi/cors v1.2.0 // indirect
	github.com/go-chi/render v1.0.1 // indirect
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/jinzhu/gorm v1.9.16
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
	github.com/travisjeffery/go-dynaport v1.0.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
)
