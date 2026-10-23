# Installing chart-verifier

The only official download source is the GitHub releases page of `redhat-certification/chart-verifier`: <https://github.com/redhat-certification/chart-verifier/releases>.
Resolve the latest release at run time unless the user requests a specific version; never pin a default version.

## Platform and prerequisites

Before running [the installer](../scripts/install-chart-verifier.sh), check the host OS and architecture with `uname -s` and `uname -m`.
Official release binaries require **Linux/amd64** (`x86_64`); the script does not check platform compatibility.

On other platforms, explain the incompatibility rather than downloading a binary that will not run.

## Download and install

Get approval for downloading, the installation destination, and running the installed binary before invoking the script.
The destination defaults to the current working directory (`./chart-verifier`):

```sh
sh <skill-dir>/scripts/install-chart-verifier.sh latest
```

Replace `latest` with an exact release tag when the user requests a specific version. Pass a second argument to install into another user-writable directory:

```sh
sh <skill-dir>/scripts/install-chart-verifier.sh latest /chosen/directory
```

The script discovers the actual tarball asset name from the official release, downloads it over HTTPS, installs the binary, and runs its version/help commands to sanity check the installation.
It refuses existing destinations, including symlinks. If installation fails, report the error rather than bypassing these safeguards.

After a successful installation, the downloaded binary should be used for further command execution.
