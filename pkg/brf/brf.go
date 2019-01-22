// brf = brief
package brf

import (
	"fmt"
	"github.com/jychri/git-in-sync/pkg/flags"
)

func Printv(f flags.Flags, s string, z ...interface{}) {

	if f.Check("oneline") {
		return
	}

	fmt.Println(fmt.Sprintf(s, z...))
}
