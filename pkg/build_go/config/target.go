package config

import (
	"fmt"

	"github.com/cardinalby/xgo-pack/pkg/consts"
)

type Target struct {
	Os   consts.Os
	Arch consts.Arch
}

func (t Target) String() string {
	return fmt.Sprintf("%s/%s", t.Os, t.Arch)
}

type Targets []Target

func (t Targets) Strings() []string {
	res := make([]string, len(t))
	for i, target := range t {
		res[i] = target.String()
	}
	return res
}
