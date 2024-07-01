//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func InstallXk6() error {
	return sh.Run(mg.GoCmd(), "install", "go.k6.io/xk6/cmd/xk6@latest")
}

func Test() error {
	mg.Deps(InstallXk6)
	return sh.Run(mg.GoCmd(), "test", "./...")
}

func Build() error {
	mg.Deps(InstallXk6)
	return sh.RunV("xk6", "build",
		"--with", "github.com/kralicky/xk6-dns=.",
		"--output", "bin/k6",
	)
}
