package modeler

import "testing"

func TestIndependentCfg(t *testing.T) {
	i := &IndependentCfg{}
	if err := i.Init(); err != nil {
		t.Error(err)
	}
}
