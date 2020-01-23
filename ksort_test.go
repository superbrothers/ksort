package ksort

import (
	"strings"
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestPrintVersionInformation(t *testing.T) {
	streams, _, _, errOut := genericclioptions.NewTestIOStreams()
	cmd := NewCommand(streams)
	cmd.SetArgs([]string{"--version"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("cmd.Execute() => %q", err)
	}

	if !strings.Contains(errOut.String(), "&ksort.info{") {
		t.Errorf("expect to include version information, but got %q", errOut)
	}
}

func TestCommand(t *testing.T) {
	tests := []struct {
		args []string
		in   string
		out  string
	}{
		{
			args: []string{"--filename", "testdata/rbac.yaml"},
			in:   "",
			out: `# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: ClusterRole
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ClusterRoleBinding
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: Role
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: RoleBinding
`,
		},
		{
			args: []string{"--filename", "testdata"},
			in:   "",
			out: `# Source: testdata/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: configmap
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: ClusterRole
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ClusterRoleBinding
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: Role
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: RoleBinding
---
# Source: testdata/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment
`,
		},
		{
			args: []string{"--filename", "testdata/deployment.yaml", "--filename", "testdata/rbac.yaml"},
			in:   "",
			out: `# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: ClusterRole
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ClusterRoleBinding
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: Role
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: RoleBinding
---
# Source: testdata/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment
`,
		},
		{
			args: []string{"--filename", "testdata/deployment.yaml", "--filename", "testdata/rbac.yaml", "--delete"},
			out: `# Source: testdata/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: RoleBinding
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: Role
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ClusterRoleBinding
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: ClusterRole
`,
		},
		{
			args: []string{"--recursive", "--filename", "testdata"},
			in:   "",
			out: `# Source: testdata/test-recursive/ns.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: test-recursive
---
# Source: testdata/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: configmap
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: ClusterRole
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ClusterRoleBinding
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: Role
---
# Source: testdata/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: RoleBinding
---
# Source: testdata/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment
`,
		},
		{
			args: []string{"--filename", "-"},
			in: `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ClusterRoleBinding
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: ClusterRole
`,
			out: `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: ClusterRole
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ClusterRoleBinding
`,
		},
	}

	for i, tt := range tests {
		streams, in, out, _ := genericclioptions.NewTestIOStreams()
		in.WriteString(tt.in)
		cmd := NewCommand(streams)
		cmd.SetArgs(tt.args)

		if err := cmd.Execute(); err != nil {
			t.Errorf("%d: cmd.Execute() => %q", i, err)
		}

		if out.String() != tt.out {
			t.Errorf("%d: out => %q, want %q", i, out.String(), tt.out)
		}
	}
}
