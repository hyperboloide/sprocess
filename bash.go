//
// shell.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Feb  8 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.
//

package sprocess

import (
	"io"
	"os/exec"
)

type Bash struct {
	Cmd string
	Name string
}

func (b *Bash) GetName() string {
	return b.Name
}

func (b *Bash) Start() error {
	_, err := exec.LookPath("bash")
	return err
}

func (b *Bash) Encode(r io.Reader, w io.Writer, d *Data) error {
	cmd := exec.Command("bash", "-lc", b.Cmd)
	cmd.Stdout = w
	cmd.Stdin = r
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}
