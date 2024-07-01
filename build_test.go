package dns_test

import (
	"fmt"
	"go/build"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestXk6Build(t *testing.T) {
	pkg, err := build.Default.Import("go.k6.io/k6", ".", build.FindOnly)
	if err != nil {
		t.Fatal(err)
	}
	_, version, _ := strings.Cut(pkg.Root, "@")

	actual, err := exec.Command("xk6", "version").Output()
	if err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf(`
k6 %s (%s, %s/%s)
Extensions:
  github.com/kralicky/xk6-dns (devel), k6/x/dns [js]

`[1:],
		version, runtime.Version(), runtime.GOOS, runtime.GOARCH)

	if string(actual) != expected {
		t.Errorf(`
expected:
---
%s
---
got:
---
%s
---
`, expected, string(actual))
	}
}
