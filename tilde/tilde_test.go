package tilde

import (
	u "os/user"
	"strings"
	"testing"
)

func TestAbsUser(t *testing.T) {

	var user *u.User
	var err error

	if user, err = u.Current(); err != nil {
		t.Errorf("Unable to identify current user")
	}

	for _, tr := range []struct {
		in, want string
	}{
		{"~/testing", strings.Join([]string{user.HomeDir, "/testing"}, "")},
		{"~/testing/", strings.Join([]string{user.HomeDir, "/testing"}, "")},
		{"/testing/", "/testing"},
	} {
		if got := Abs(tr.in); err != nil || got != tr.want {
			t.Errorf("AbsUser: (%v != %v)", got, tr.want)
		}
	}
}
