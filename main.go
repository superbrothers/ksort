/*
Copyright (c) 2019 Kazuki Suda <kazuki.suda@gmail.com>

For the full copyright and license information, please view the LICENSE
file that was distributed with this source code.
*/

package main

import (
	"os"

	"github.com/superbrothers/ksort/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
