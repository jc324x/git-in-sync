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

func RemoveDuplicates(ssl []string) (sl []string) {

	smap := make(map[string]bool)

	for i := range ssl {
		if smap[ssl[i]] == true {
		} else {
			smap[ssl[i]] = true
			sl = append(sl, ssl[i])
		}
	}

	return sl
}
