package stat

import (
	"testing"
)

func TestInit(t *testing.T) {
	st := Init()
	st.Reduce()
}
