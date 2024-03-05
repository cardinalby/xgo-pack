package plist

import v "github.com/go-ozzo/ozzo-validation/v4"

type Properties struct {
	ProductName             string
	BinName                 string
	Identifier              string
	ProductVersion          string
	IconFile                string
	IsHighResolutionCapable bool
	IsHideInDock            bool
	CopyRight               string
}

func (p *Properties) Validate() error {
	return v.ValidateStruct(p,
		v.Field(&p.BinName, v.Required),
		v.Field(&p.Identifier, v.Required),
	)
}
