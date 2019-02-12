package tilde

import (
	"os/user"
	"strings"
	"testing"
)

func TestAbsUser(t *testing.T) {

	u, err := user.Current()

	if err != nil {
		t.Errorf("Relative: Can't identify current user")
	}

	for _, tr := range []struct {
		in, want string
	}{
		{"~/testing", strings.Join([]string{u.HomeDir, "/testing"}, "")},
		{"~/testing/", strings.Join([]string{u.HomeDir, "/testing"}, "")},
		{"/testing/", "/testing"},
	} {
		if got := Abs(tr.in); err != nil || got != tr.want {
			t.Errorf("AbsUser: (%v != %v)", got, tr.want)
		}
	}
}
