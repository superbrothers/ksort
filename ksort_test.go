package ksort

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestCommand(t *testing.T) {
	tests := []struct {
		args []string
		out  string
	}{
		{
			args: []string{"testdata/rbac.yaml"},
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
			args: []string{"testdata"},
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
			args: []string{"testdata/deployment.yaml", "testdata/rbac.yaml"},
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
			args: []string{"testdata/deployment.yaml", "testdata/rbac.yaml", "--delete"},
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
	}

	for i, tt := range tests {
		out := &bytes.Buffer{}
		cmd := NewCommand(bytes.NewReader(nil), out, ioutil.Discard)
		cmd.SetArgs(tt.args)

		if err := cmd.Execute(); err != nil {
			t.Errorf("%d: cmd.Execute() => %q", i, err)
		}

		if out.String() != tt.out {
			t.Errorf("%d: out => %q, want %q", i, out.String(), tt.out)
		}
	}
}
