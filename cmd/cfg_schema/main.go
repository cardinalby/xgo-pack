package main

import (
	"encoding/json"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/cfgtypes"
	fsutil "github.com/cardinalby/xgo-pack/pkg/util/fs"
	"github.com/invopop/jsonschema"
)

func main() {
	r := new(jsonschema.Reflector)
	err := r.AddGoComments("github.com/cardinalby/xgo-pack", "pkg")
	if err != nil {
		panic(err)
	}
	r.ExpandedStruct = true
	s := r.Reflect(&cfgtypes.Config{})
	s.ID = cfgtypes.DefaultSchema

	data, err := json.MarshalIndent(s, "", "  ")

	if err != nil {
		panic(err)
	}
	if err := fsutil.WriteFile("config_schema/config.schema.v1.json", data); err != nil {
		panic(err)
	}
}
