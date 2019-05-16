package ksort

import (
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestCommand(t *testing.T) {
	tests := []struct {
		args []string
		out  string
	}{
		{
			args: []string{"--filename", "testdata/rbac.yaml"},
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
	}

	for i, tt := range tests {
		streams, _, out, _ := genericclioptions.NewTestIOStreams()
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
