# OpenShift Version in the Report

The chart-verifier tool adds the OpenShift version where chart-testing has run
in the report.  The information is available as a metadata at this path:
`.metadata.tool.certifiedOpenShiftVersions`.  Here is an example showing only
the relevant details:

```
...
metadata:
  tool:
    certifiedOpenShiftVersions: 4.7.8
```

The chart-verifier tool runs the `oc version -o yaml` command to retrieve the
OpenShift version value.  It gives the version value if the logged-in user
(role) has access to `get` values of `clusteroperators` (a cluster scoped
resource in the `config.openshift.io` API group).  You need to configure a
specific role for the user, as given here.

You need a ClusterRole like this:

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

And ClusterRoleBinding like this:

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

You can replace the name of the ServiceAccount and its namespace values in the
above example. Note: Instead of a service account, the [subject could be a user
or group][subject].

If giving the role mentioned above is not feasible, alternatively, you can
specify the OpenShift version as a command-line flag in the chart-verifier
tool. You can use the `--openshift-version` flag to specify the version. Here is
an example:

```
$ podman run -it --rm quay.io/redhat-certification/chart-verifier verify --openshift-version=4.9.2 <chart-uri>
```

The `-V` flag is the short version for the `--openshift-version` flag.

```
$ podman run -it --rm quay.io/redhat-certification/chart-verifier verify -V 4.9.2 <chart-uri>
```

[subject]: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#referring-to-subjects
