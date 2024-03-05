package manifest

import v "github.com/go-ozzo/ozzo-validation/v4"

type Properties struct {
	Identifier string
	IsDpiAware bool
}

func (p *Properties) Validate() error {
	return v.ValidateStruct(p,
		v.Field(&p.Identifier, v.Required),
	)
}
