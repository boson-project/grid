package platformb

import (
	"github.com/boson-project/grid/mock"
)

type Adapter struct {
	mock.Adapter
}

func NewAdapter() Adapter {
	return Adapter{}
}
