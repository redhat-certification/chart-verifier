# OpenShift Version in the Report

The chart-verifier CLI tool adds the OpenShift version where chart-testing has run
in the report.  The information is available as a metadata at this path:
`.metadata.tool.certifiedOpenShiftVersions`.

**Example**

```
...
metadata:
  tool:
    certifiedOpenShiftVersions: 4.7.8
```

The chart-verifier CLI tool runs the `oc version -o yaml` command to retrieve
the value of the OpenShift version.  The `oc version` command gives the version
value if the logged-in user (role) has access to `get` values of
`clusteroperators` (a cluster scoped resource in the `config.openshift.io` API
group).  You need to configure a specific role (ClusterRole) for the user, as
shown in the following example:

**Example of a ClusterRole**

```
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: chart-verifier-cluster-operator-role
rules:
  - apiGroups:
      - "config.openshift.io"
    resources:
      - 'clusteroperators'
    verbs:
      - 'get'
```

**Example of a ClusterRoleBinding**

```
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: chart-verifier-cluster-operator-role-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: chart-verifier-cluster-operator-role
subjects:
  - kind: ServiceAccount
    name: <SERVICE-ACCOUNT-NAME>
    namespace: <SERVICE-ACCOUNT-NAMESPACE>
```

**Notes:**

1. You can replace values of the name and namespace of the ServiceAccount in the
   previous example.
2. Instead of a service account, the [subject could be a user or
   group][subject].

 If configuring the role as mentioned in the previous examples is not feasible,
alternatively, you can specify the OpenShift version as a command-line flag in
the chart-verifier CLI tool.

For example, you can use the `--openshift-version` flag to specify the OpenShift
version in the chart-verifier CLI tool:

```
$ podman run -it --rm quay.io/redhat-certification/chart-verifier verify --openshift-version=4.9.2 <chart-uri>
```

Note: The `-V` flag is the short version for the `--openshift-version` flag.

**Example**

```
$ podman run -it --rm quay.io/redhat-certification/chart-verifier verify -V 4.9.2 <chart-uri>
```

[subject]: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#referring-to-subjects
