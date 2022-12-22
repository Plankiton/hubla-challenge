package main

type Validator struct {
}

func NewValidator() (*Validator, error) {
	err := resterrors.InitValidator()
	if err != nil {
		return nil, err
	}

	return &Validator{}, nil
}

func (v *Validator) Validate(i interface{}) error {
	return resterrors.Validator().Struct(i)
}
