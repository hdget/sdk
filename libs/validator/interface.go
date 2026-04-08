package validator

type Validator interface {
	Validate(input any) error
	MustValidate(input any)
}
