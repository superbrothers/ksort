/*
Copyright (c) 2019 Kazuki Suda <kazuki.suda@gmail.com>

For the full copyright and license information, please view the LICENSE
file that was distributed with this source code.
*/

package ksort

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	util "k8s.io/helm/pkg/releaseutil"
	"k8s.io/helm/pkg/tiller"
)

var (
	printVersion = false
)

const (
	ksortLong = `When installing manifests, they should be sorted in a proper order by Kind.
For example, Namespace object must be in the first place when installing them.

ksort sorts manfest files in a proper order by Kind, which is implementd by
using SortByKind function in Kubernetes Helm.`

	ksortExample = `  # Sort manifest files under the manifests directory, and output the result to the stdout.
  ksort ./manifests

  # To pass the result into the stdin of kubectl apply command is also convenient.
  ksort ./manifests | kubectl apply -f -

  # Sort manifests contained the manifest file.
  ksort app.yaml`

	kindUnknown = "Unknown"
)

type options struct {
	filenames []string
}

func init() {
	flag.Set("logtostderr", "true")
}

func NewCommand(in io.Reader, out, errOut io.Writer) *cobra.Command {
	o := options{}

	cmd := &cobra.Command{
		Use:     "ksort FILENAME...",
		Short:   "ksort sorts manfest files in a proper order by Kind.",
		Long:    ksortLong,
		Example: ksortExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if printVersion {
				fmt.Fprintf(errOut, "%#v\n", newInfo())
				return nil
			}

			if err := o.complete(cmd, args); err != nil {
				return err
			}

			cmd.SilenceUsage = true

			return o.run(out)
		},
	}

	cmd.SetOutput(errOut)

	cmd.Flags().BoolVar(&printVersion, "version", printVersion, "Print the version and exit")
	cmd.Flags().AddGoFlagSet(flag.CommandLine)

	// Workaround for this issue:
	// https://github.com/kubernetes/kubernetes/issues/17162
	flag.CommandLine.Parse([]string{})

	return cmd
}

func (o *options) complete(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("filename is required")
	}

	// verify manifest file exists
	for _, filename := range args {
		if _, err := os.Stat(filename); err != nil {
			return err
		}
	}

	o.filenames = args

	return nil
}

func (o *options) run(out io.Writer) error {
	contents := map[string]string{}

	for _, filename := range o.filenames {
		glog.V(2).Infof("Walking the file tree rooted at %q", filename)

		err := filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
			glog.V(2).Infof("Visiting %q", path)

			if err != nil {
				return fmt.Errorf("Failed to access a path %q: %v\n", filename, err)
			}

			if info.IsDir() {
				glog.V(2).Infof("Skip %q because it's directory", path)
				return nil
			}

			if _, ok := contents[path]; ok {
				glog.V(2).Infof("Skip reading %q because it already went through", path)
				return nil
			}

			content, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("Failed to read a file %q: %v\n", path, err)
			}

			contents[path] = string(content)

			return nil
		})

		if err != nil {
			return err
		}
	}

	if len(contents) == 0 {
		return errors.New("File does not exist")
	}

	// extract kind and name
	re := regexp.MustCompile("kind:(.*)\n")
	// YAML separator
	sep := regexp.MustCompile("(?m)^---.*$")
	manifests := []tiller.Manifest{}
	for k, v := range contents {
		docs := sep.Split(v, -1)
		for _, doc := range docs {
			if len(doc) == 0 {
				continue
			}
			match := re.FindStringSubmatch(doc)
			h := kindUnknown
			if len(match) == 2 {
				h = strings.TrimSpace(match[1])
			}
			doc = strings.Trim(doc, "\n")
			m := tiller.Manifest{Name: k, Content: doc, Head: &util.SimpleHead{Kind: h}}
			manifests = append(manifests, m)
			glog.V(2).Infof("Found %s in %q", h, k)
		}
	}
	glog.V(2).Infof("Found %d objects in total", len(manifests))

	a := make([]string, len(manifests))
	for i, m := range tiller.SortByKind(manifests) {
		a[i] += fmt.Sprintf("# Source: %s\n", m.Name)

		if m.Head.Kind == kindUnknown {
			a[i] += "# WARNING: It looks like that this file is not a manifest file\n"
			continue
		}

		a[i] += m.Content
	}

	fmt.Fprintln(out, strings.Join(a, "\n---\n"))

	return nil
}
