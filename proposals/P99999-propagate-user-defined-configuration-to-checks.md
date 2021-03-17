# Propagate User Defined Configuration To Checks

* Proposal N.: 99999
* Authors: isuttonl@redhat.com
* Status: **Draft**

## Abstract

There are multiple use cases where propagating user defined configuration to checks is useful: to specify an OpenShift
version to verify the chart compatibility or influence the severity the Helm linter check should consider a failure; in
other words, any occasion the user needs to parametrize a check.

## Rationale

`chart-verifier` is the command (either through a Container Runtime or directly using the program directly) users
interface to provide information regarding the verifications to be performed for a given chart, so it is expected that
different parameters can be used at different moments in time.

Those parameters can be given to the program in two ways: through the configuration file, which is already supported by
not yet used; and through command line flags that can be used to overwrite a value defined by the configuration file.

Having both mechanisms to influence the verification session is very useful for a couple of reasons:

1. Configuration files can be distributed as *profiles*, for example one for OpenShift 4.6, another for OpenShift 4.7;
1. Profile defaults can be overridden, helping developers and other power users to debug, change checks parameters
   before committing to write a configuration/profile file.

## Usage

Some changes are required in the `chart-verifier` program, more specifically in the `verify` command, where the
flag `--set` should be introduced, as in the example below:

```text
> chart-verifier verify --set compat.version=openshift-4.6 --enable compat chart.tgz
```

As the example above shows, the format of a `--set` flag is `KEY=VALUE`, where `KEY` is the path of the value in the
configuration file and `VALUE` the value to be assigned to the configuration.

In the example, `compat` is the check name, and `version` is the key the `compat` check will use to verify the
compatibility of the given chart, in this case against the `openshift-4.6` profile.

Another example related to the severity of the Helm linter check:

```text
> chart-verifier verify --set linter.failWhen=ERROR --enable linter chart.tgz
```

Both configurations could be use simultaneously by providing multiple `--set` flags:

```text
> chart-verifier verify                \
    --set compat.version=openshift-4.6 \
    --set linter.failWhen=ERROR        \
    --enable compat,linter             \
    chart.tgz
```

The configuration file could be used to store the same configuration expressed by the usage of the `--set` flag above:

```yaml
compat:
    version: openshift-4.6
linter:
    failWhen: ERROR
```

The example below uses the configuration file above:

```text
> chart-verifier verify --config openshift-4.6.yaml chart.tgz
```

## Checks API

Currently the `CheckFunc` type is defined as the example below:

```go
type CheckFunc func(uri string) (Result, error)
```

The current `CheckFunc` type is too simplistic, and should be modified to either receive an extra parameter containing
the options to be considered by the check, or a more complex type composing both of this values:

```go
type CheckFunc func(uri string, config *viper.Viper) (Result, error)
```

`*viper.Viper` is used since it has already been integrated in the main program, and also because it offers the required configuration semantics to isolate checks configuration.

The `opts` map is configuration for the check found in the configuration file overridden by the `--set`; in the `openshift-4.6.yaml`, the `compat` check would expect the following values in `opts`:

```go
opts := map[string]interface{}{
    "version": "openshift-4.6"
}
```

## Test Scenarios

There are a couple of scenarios that should be tested, in order to validate the proposal implementation (assuming a `dummy` check is available):

* The `dummy` check should fail by default, regardless the chart;
* The `dummy` check should succeed when `dummy.ok` is set to `true` using the `--set` flag;
* The `dummy` check should succeed when the configuration file configures the `dummy.ok` configuration value to `true`; and
* The `dummy` check should fail when the configuration file configures the `dummy.ok` configuration value to `true` and `dummy.ok` is set to `false` using the `--set` flag.