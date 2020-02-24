package util

import (
	"os"
	"os/exec"
	"strings"
)

type Pager struct {
	strings.Builder
}

func NewPager() *Pager {
	return new(Pager)
}

func (p *Pager) Show() {

	cmd := exec.Command("less", "-FRX")
	cmd.Stdin = strings.NewReader(p.Builder.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
