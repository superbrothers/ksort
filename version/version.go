/*
Copyright (c) 2019 Kazuki Suda <kazuki.suda@gmail.com>

For the full copyright and license information, please view the LICENSE
file that was distributed with this source code.
*/

package version

import (
	"fmt"
	"runtime"
)

var (
	GitVersion   = "v0.0.0"
	GitCommit    = "$Format:%H$"
	GitTreeState = "dirty"
	BuildDate    = "1970-01-01T00:00:00Z"
)

type Info struct {
	GitVersion   string
	GitCommit    string
	GitTreeState string
	BuildDate    string
	GoVersion    string
	Compiler     string
	Platform     string
}

func NewInfo() *Info {
	return &Info{
		GitVersion:   GitVersion,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
