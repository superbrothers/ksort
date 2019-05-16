/*
Copyright (c) 2019 Kazuki Suda <kazuki.suda@gmail.com>

For the full copyright and license information, please view the LICENSE
file that was distributed with this source code.
*/

package main

import (
	"flag"
	"os"

	"github.com/superbrothers/ksort"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func init() {
	flag.Set("logtostderr", "true")
}

func main() {
	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	cmd := ksort.NewCommand(streams)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
