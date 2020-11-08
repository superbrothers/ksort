/*
Copyright (c) 2019 Kazuki Suda <kazuki.suda@gmail.com>

For the full copyright and license information, please view the LICENSE
file that was distributed with this source code.
*/

package ksort

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/releaseutil"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/klog/v2"
)

const (
	ksortLong = `When installing manifests, they should be sorted in a proper order by Kind.
For example, Namespace object must be in the first place when installing them.

ksort sorts manfest files in a proper order by Kind, which is implementd by
using sortManifestsByKind function in Kubernetes Helm.`

	ksortExample = `  # Sort manifest files under the manifests directory, and output the result to the stdout.
  ksort -f ./manifests

  # To pass the result into the stdin of kubectl apply command is also convenient.
  ksort -f ./manifests | kubectl apply -f -

  # Sort manifests contained the manifest file.
  ksort -f app.yaml

  # Sort manifests in uninstall order.
  ksort -f ./manifests --delete

  # Sort manifests passed into stdin.
  cat app.yaml | ksort -f -`

	kindUnknown = "Unknown"
)

type options struct {
	filenameFlags *genericclioptions.FileNameFlags

	filenameOptions resource.FilenameOptions
	delete          bool

	genericclioptions.IOStreams
}

func newOptions(streams genericclioptions.IOStreams) *options {
	usage := "Containing the resource to sort in a proper order by Kind"

	filenames := []string{}
	recursive := false

	return &options{
		filenameFlags: &genericclioptions.FileNameFlags{
			Usage:     usage,
			Filenames: &filenames,
			Recursive: &recursive,
		},
		IOStreams: streams,
	}
}

func NewCommand(streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(streams)

	printVersion := false

	cmd := &cobra.Command{
		Use:     "ksort -f FILENAME",
		Short:   "ksort sorts manfest files in a proper order by Kind.",
		Long:    ksortLong,
		Example: ksortExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if printVersion {
				fmt.Fprintf(o.ErrOut, "%#v\n", newInfo())
				return nil
			}

			if err := o.complete(cmd, args); err != nil {
				return err
			}

			if err := o.validate(); err != nil {
				return err
			}

			cmd.SilenceUsage = true

			return o.run()
		},
	}

	o.filenameFlags.AddFlags(cmd.Flags())
	cmd.Flags().BoolVarP(&o.delete, "delete", "d", o.delete, "Sort manifests in uninstall order")
	cmd.Flags().BoolVar(&printVersion, "version", printVersion, "Print the version and exit")
	cmd.Flags().AddGoFlagSet(flag.CommandLine)

	// Workaround for this issue:
	// https://github.com/kubernetes/kubernetes/issues/17162
	_ = flag.CommandLine.Parse([]string{})

	return cmd
}

func (o *options) complete(cmd *cobra.Command, args []string) error {
	o.filenameOptions = o.filenameFlags.ToOptions()

	return nil
}

func (o *options) validate() error {
	if len(o.filenameOptions.Filenames) == 0 {
		return errors.New("Must specify --filename")
	}

	return nil
}

func (o *options) run() error {
	contents := map[string]string{}

	for _, filename := range o.filenameOptions.Filenames {
		if filename == "-" {
			klog.V(2).Infof("Reading manifest from the standard input")

			var lines []string
			scanner := bufio.NewScanner(o.In)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			contents[""] = strings.Join(lines, "\n")
			continue
		}

		klog.V(2).Infof("Walking the file tree rooted at %q", filename)

		err := filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
			klog.V(2).Infof("Visiting %q", path)

			if err != nil {
				return fmt.Errorf("Failed to access a path %q: %v\n", filename, err)
			}

			if info.IsDir() {
				if path != filename && !o.filenameOptions.Recursive {
					return filepath.SkipDir
				}
				return nil
			}

			if _, ok := contents[path]; ok {
				klog.V(2).Infof("Skip reading %q because it already went through", path)
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
		return errors.New("There are no files")
	}

	// extract kind and name
	re := regexp.MustCompile("kind:(.*)\n")
	// YAML separator
	sep := regexp.MustCompile("(?m)^---.*$")
	manifests := []releaseutil.Manifest{}
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
			m := releaseutil.Manifest{Name: k, Content: doc, Head: &releaseutil.SimpleHead{Kind: h}}
			manifests = append(manifests, m)
			klog.V(2).Infof("Found %s in %q", h, k)
		}
	}
	klog.V(2).Infof("Found %d objects in total", len(manifests))

	var sortOrder KindSortOrder
	if o.delete {
		sortOrder = UninstallOrder
	} else {
		sortOrder = InstallOrder
	}

	a := make([]string, len(manifests))
	for i, m := range sortManifestsByKind(manifests, sortOrder) {
		// If manifest data is read from stdin, m.Name is empty
		if m.Name != "" {
			a[i] += fmt.Sprintf("# Source: %s\n", m.Name)
		}

		if m.Head.Kind == kindUnknown {
			a[i] += "# WARNING: It looks like that this file is not a manifest file\n"
			continue
		}

		a[i] += m.Content
	}

	fmt.Fprintln(o.Out, strings.Join(a, "\n---\n"))

	return nil
}
