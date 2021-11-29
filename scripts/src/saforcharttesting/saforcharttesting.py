import sys
import time
import os
import base64
import json
import argparse
import subprocess
import tempfile
from string import Template

namespace_template = """\
apiVersion: v1
kind: Namespace
metadata:
  name: ${name}
"""

serviceaccount_template = """\
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ${name}
  namespace: ${name}
"""

role_template = """\
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: ${name}
  namespace: ${name}
rules:
  - apiGroups:
      - "*"
    resources:
      - '*'
    verbs:
      - '*'
"""

rolebinding_template = """\
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: ${name}
  namespace: ${name}
subjects:
- kind: ServiceAccount
  name: ${name}
  namespace: ${name}
roleRef:
  kind: Role
  name: ${name}
"""

clusterrole_template = """\
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: ${name}
rules:
  - apiGroups:
      - "config.openshift.io"
    resources:
      - 'clusteroperators'
    verbs:
      - 'get'
  - apiGroups:
      - "rbac.authorization.k8s.io"
    resources:
      - 'clusterrolebindings'
      - 'clusterroles'
    verbs:
      - 'get'
      - 'create'
      - 'delete'
  - apiGroups:
      - "admissionregistration.k8s.io"
    resources:
      - 'mutatingwebhookconfigurations'
    verbs:
      - 'get'
      - 'create'
      - 'list'
      - 'watch'
      - 'patch'
      - 'delete'
  - apiGroups:
      - "authentication.k8s.io"
    resources:
      - 'tokenreviews'
    verbs:
      - 'create'
  - apiGroups:
      - "authorization.k8s.io"
    resources:
      - 'subjectaccessreviews'
    verbs:
      - 'create'
"""

clusterrolebinding_template = """\
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: ${name}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: ${name}
subjects:
  - kind: ServiceAccount
    name: ${name}
    namespace: ${name}
"""

def apply_config(tmpl, **values):
    with tempfile.TemporaryDirectory(prefix="sa-for-chart-testing-") as tmpdir:
        content = Template(tmpl).substitute(values)
        config_path = os.path.join(tmpdir, "config.yaml")
        with open(config_path, "w") as fd:
            fd.write(content)
        out = subprocess.run(["./oc", "apply", "-f", config_path], capture_output=True)
        stdout = out.stdout.decode("utf-8")
        if out.returncode != 0:
            stderr = out.stderr.decode("utf-8")
        else:
            stderr = ""

    return stdout, stderr

def delete_config(tmpl, **values):
    with tempfile.TemporaryDirectory(prefix="sa-for-chart-testing-") as tmpdir:
        content = Template(tmpl).substitute(values)
        config_path = os.path.join(tmpdir, "config.yaml")
        with open(config_path, "w") as fd:
            fd.write(content)
        out = subprocess.run(["./oc", "delete", "-f", config_path], capture_output=True)
        stdout = out.stdout.decode("utf-8")
        if out.returncode != 0:
            stderr = out.stderr.decode("utf-8")
        else:
            stderr = ""

    return stdout, stderr

def create_namespace(namespace):
    print("creating Namespace:", namespace)
    stdout, stderr = apply_config(namespace_template, name=namespace)
    print("stdout:\n", stdout, sep="")
    if stderr.strip():
        print("[ERROR] creating Namespace:", stderr)

def create_serviceaccount(namespace):
    print("creating ServiceAccount:", namespace)
    stdout, stderr = apply_config(serviceaccount_template, name=namespace)
    print("stdout:\n", stdout, sep="")
    if stderr.strip():
        print("[ERROR] creating ServiceAccount:", stderr)

def create_role(namespace):
    print("creating Role:", namespace)
    stdout, stderr = apply_config(role_template, name=namespace)
    print("stdout:\n", stdout, sep="")
    if stderr.strip():
        print("[ERROR] creating Role:", stderr)

def create_rolebinding(namespace):
    print("creating RoleBinding:", namespace)
    stdout, stderr = apply_config(rolebinding_template, name=namespace)
    print("stdout:\n", stdout, sep="")
    if stderr.strip():
        print("[ERROR] creating RoleBinding:", stderr)

def create_clusterrole(namespace):
    print("creating ClusterRole:", namespace)
    stdout, stderr = apply_config(clusterrole_template, name=namespace)
    print("stdout:\n", stdout, sep="")
    if stderr.strip():
        print("[ERROR] creating ClusterRole:", stderr)

def create_clusterrolebinding(namespace):
    print("creating ClusterRoleBinding:", namespace)
    stdout, stderr = apply_config(clusterrolebinding_template, name=namespace)
    print("stdout:\n", stdout, sep="")
    if stderr.strip():
        print("[ERROR] creating ClusterRoleBinding:", stderr)

def delete_namespace(namespace):
    print("deleting Namespace:", namespace)
    stdout, stderr = delete_config(namespace_template, name=namespace)
    print("stdout:\n", stdout, sep="")
    if stderr.strip():
        print("[ERROR] deleting Namespace:", namespace, stderr)
        sys.exit(1)

def delete_clusterrole(name):
    print("deleting ClusterRole:", name)
    stdout, stderr = delete_config(clusterrole_template, name=name)
    print("stdout:\n", stdout, sep="")
    if stderr.strip():
        print("[ERROR] deleting ClusterRole:", name, stderr)
        sys.exit(1)

def delete_clusterrolebinding(name):
    print("deleting ClusterRoleBinding:", name)
    stdout, stderr = delete_config(clusterrolebinding_template, name=name)
    print("stdout:\n", stdout, sep="")
    if stderr.strip():
        print("[ERROR] deleting ClusterRoleBinding:", name, stderr)
        sys.exit(1)

def write_sa_token(namespace, token):
    sa_found = False
    for i in range(7):
        out = subprocess.run(["./oc", "get", "serviceaccount", namespace, "-n", namespace, "-o", "json"], capture_output=True)
        stdout = out.stdout.decode("utf-8")
        if out.returncode != 0:
            stderr = out.stderr.decode("utf-8")
            if stderr.strip():
                print("[ERROR] retrieving ServiceAccount:", namespace, stderr)
                time.sleep(10)
        else:
            sa = json.loads(stdout)
            if len(sa["secrets"]) >= 2:
                sa_found = True
                break
            time.sleep(10)

    if not sa_found:
        print("[ERROR] retrieving ServiceAccount:", namespace, stderr)
        sys.exit(1)

    for secret in sa["secrets"]:
        out = subprocess.run(["./oc", "get", "secret", secret["name"], "-n", namespace, "-o", "json"], capture_output=True)
        stdout = out.stdout.decode("utf-8")
        if out.returncode != 0:
            stderr = out.stderr.decode("utf-8")
            if stderr.strip():
                print("[ERROR] retrieving secret:", secret["name"], stderr)
                continue
        else:
            sec = json.loads(stdout)
            if sec["type"] == "kubernetes.io/service-account-token":
                content = sec["data"]["token"]
                with open(token, "w") as fd:
                    fd.write(base64.b64decode(content).decode("utf-8"))

def switch_project_context(namespace, token, api_server):
    tkn = open(token).read()
    for i in range(7):
        out = subprocess.run(["./oc", "login", "--token", tkn, "--server", api_server], capture_output=True)
        stdout = out.stdout.decode("utf-8")
        print(stdout)
        out = subprocess.run(["./oc", "project", namespace], capture_output=True)
        stdout = out.stdout.decode("utf-8")
        print(stdout)
        out = subprocess.run(["./oc", "config", "current-context"], capture_output=True)
        stdout = out.stdout.decode("utf-8").strip()
        print(stdout)
        if stdout.endswith(":".join((namespace, namespace))):
            print("current-context:", stdout)
            return
        time.sleep(10)

    # This exit will happen if there is an infra failure
    print("""[ERROR] There is an error creating the namespace and service account. It happens due to some infrastructure failure.  It is not directly related to the changes in the pull request. You can wait for some time and try to re-run the job.  To re-run the job change the PR into a draft and remove the draft state.""")
    sys.exit(1)

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--create", dest="create", type=str, required=False,
                                        help="create service account and namespace for chart testing")
    parser.add_argument("-t", "--token", dest="token", type=str, required=False,
                                        help="service account token for chart testing")
    parser.add_argument("-d", "--delete", dest="delete", type=str, required=False,
                                        help="delete service account and namespace used for chart testing")
    parser.add_argument("-s", "--server", dest="server", type=str, required=False,
                                        help="API server URL")
    args = parser.parse_args()

    if args.create:
        create_namespace(args.create)
        create_serviceaccount(args.create)
        create_role(args.create)
        create_rolebinding(args.create)
        create_clusterrole(args.create)
        create_clusterrolebinding(args.create)
        write_sa_token(args.create, args.token)
        switch_project_context(args.create, args.token, args.server)
    elif args.delete:
        delete_clusterrolebinding(args.delete)
        delete_clusterrole(args.delete)
        delete_namespace(args.delete)
    else:
        parser.print_help()
