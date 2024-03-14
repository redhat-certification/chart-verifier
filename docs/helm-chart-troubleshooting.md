# Troubleshooting

- [Troubleshooting](#troubleshooting)
  - [Troubleshooting check failures](#troubleshooting-check-failures)
    - [`is-helm-v3` v1.0](#is-helm-v3-v10)
    - [`has-readme` v1.0](#has-readme-v10)
    - [`contains-test` v1.0](#contains-test-v10)
    - [`has-kubeversion` v1.0](#has-kubeversion-v10)
    - [`has-kubeversion` v1.1](#has-kubeversion-v11)
    - [`contains-values` v1.0](#contains-values-v10)
    - [`contains-values-schema` v1.0](#contains-values-schema-v10)
    - [`not-contains-crds` v1.0](#not-contains-crds-v10)
    - [`not-contain-csi-objects` v1.0](#not-contain-csi-objects-v10)
    - [`helm-lint` v1.0](#helm-lint-v10)
    - [`images-are-certified` v1.0](#images-are-certified-v10)
    - [`images-are-certified` v1.1](#images-are-certified-v11)
    - [`chart-testing` v1.0](#chart-testing-v10)
    - [`required-annotations-present` v1.0](#required-annotations-present-v10)
    - [`signature-is-valid` v1.0](#signature-is-valid-v10)
  - [Report related submission failures](#report-related-submission-failures)
    - [One or more mandatory checks have failed or are missing from the report.](#one-or-more-mandatory-checks-have-failed-or-are-missing-from-the-report)
    - [The digest in the report does not match the digest calculated for the submitted chart.](#the-digest-in-the-report-does-not-match-the-digest-calculated-for-the-submitted-chart)
    - [The certifiedOpenShiftVersions annotation does not contain a valid value.](#the-certifiedopenshiftversions-annotation-does-not-contain-a-valid-value)
  - [The chart uri is not a valid url.](#the-chart-uri-is-not-a-valid-url)
    
## Troubleshooting check failures

### `is-helm-v3` v1.0

Requires the "api-version" attribute of chart.yaml to be set to "v2". Any other value will result in the check failing.

### `has-readme` v1.0

Requires a "README.md" file to exist in the root directory of the chart. Any other spelling or
capitialisation of letters will result in the check failing.

### `contains-test` v1.0

Requires at least one file to exist in the ```templates/tests``` subdirectory of the chart. If no such file
exists this check will fail. Note the `chart-testing` check will require the directory to contain a valid test.

See also helm documentation: [chart tests](https://helm.sh/docs/topics/chart_tests/)

### `has-kubeversion` v1.0

Requires the "kubeVersion" attribute of chart.yaml to be set to a value. If the attribute is not set the check
will fail. The value set is not checked.

### `has-kubeversion` v1.1 

Requires the ```kubeVersion``` attribute of chart.yaml to be set to a valid semantic version. If the attribute is not a valid semantic version the check will fail.

See also helm documentation: [Helm documentetaion of the kubeVersion attribute](https://helm.sh/docs/topics/charts/#the-kubeversion-field)

Note: The kubeVersion field will be used to detremine the OpenShift versions the charts supports and will be set as annotation.

### `contains-values` v1.0

Requires a ```values.yaml``` file to be present in the chart. If the file is not present the check will fail.

See also helm documentation: [values](https://helm.sh/docs/chart_template_guide/values_files/) and [Best Practices for using values](https://helm.sh/docs/chart_best_practices/values/).

### `contains-values-schema` v1.0

Requires a ```values.schema.json``` file to be present in the chart. If the file is not present the check will fail.

See also helm documentation: [Schema Files](https://helm.sh/docs/topics/charts/#schema-files)

### `not-contains-crds` v1.0

Requires no RCRD's to be defined in the chart. A crd is a file with an extension of `.yaml`, `.yml` or `.json`
in a `crd` subdirectory of the chart and should be removed if present.

CRD's should be defined using operators. See: [Operator CRDs](https://docs.openshift.com/container-platform/4.2/operators/crds/crd-extending-api-with-crds.html)

### `not-contain-csi-objects` v1.0

Requires no csi objects in a chart. A csi object is a file in the template subdirectory, with an extension of `.yaml`,
and containing an `kind` attribute set to `CSIDriver`. If such a file exists it should be removed.


### `helm-lint` v1.0

Requires a `helm lint` of the chart to not result in any `ERROR` messages. If an ERROR does occur the helm lint messages
will be output. Run `helm lint` on your chart for additional information. If the chart requires specification of additional
values to pass `helm lint` use one of the `chart-set` flags of the verifier tool for this check to pass. If additional
values are required a verifier report mut be included in the chart submission.

### `images-are-certified` v1.0

Requires any images referenced in a chart to be Red Hat Certified.
- The list of image references is found by running `helm template` and if this fails the error output from `helm template`
  will be output. Run `helm template` on your chart for additional information. If the chart requires specification of additional
  attributes to pass `helm template` use one of the `chart-set` flags of the verifier tool for this check to pass. If additional
  attributes are required a verifier report must be included in the chart submission.
- Each image reference found from helm template is parsed to determine the registry, repository and tag or digest value.
    - registry is the string before the first "/" in the image reference but only if it includes a "." character.
    - the repository is what remains in the image reference, after the registry is removed and before ":" or "@sha"
    - tag is what is set after the ":" character
    - digest is what is set after the "@" character in "@sha"
- If a registry is not found the pyxis swagger api is used to find the repository and from it, extract the registry
    - `https://catalog.redhat.com/api/containers/v1/repositories?filter=repository==<repository>`
    - if the repository is not found the check will fail.
- The registry and repository are then used to find images:
    - `https://catalog.redhat.com/api/containers/v1/repositories/registry/<registry>/repository/<repository>/images`
    - if the image specified a sha value it is compared with the `parsed_data.docker_image_digest` attribute. If a
      match is not found the check fails.
    - if the image specified a tag value it is compared with the `repositories.tags.name` attributes. If a match is
      not found the check fails.
- If the check fails use the point of failure to determine how to address the issue. 

### `images-are-certified` v1.1

Requires any images referenced in a chart to be Red Hat Certified.
- The list of image references is found by running `helm template` and if this fails the error output from `helm template`
  will be output. Run `helm template` on your chart for additional information. If the chart requires specification of additional
  attributes to pass `helm template` use one of the `chart-set` flags of the verifier tool for this check to pass. If additional
  attributes are required a verifier report must be included in the chart submission.
- Each image reference found from helm template is parsed to determine the registry, repository and tag or digest value.
    - registry is the string before the first "/" in the image reference but only if it includes a "." character.
    - the repository is what remains in the image reference, after the registry is removed and before ":" or "@sha"
    - tag is what is set after the ":" character
    - digest is what is set after the "@" character in "@sha"
- If a registry is not found the pyxis swagger api is used to find the repository and from it, extract the registry
    - `https://catalog.redhat.com/api/containers/v1/repositories?filter=repository==<repository>`
    - if the repository is not found the check will fail.
- The registry and repository are then used to find images:
    - `https://catalog.redhat.com/api/containers/v1/repositories/registry/<registry>/repository/<repository>/images`
    - if the image specified a sha value it is compared with the `parsed_data.docker_image_digest` attribute. If a
      match is not found and the registry is not registry.redhat.io, the check fails.
    - if the image specified a tag value it is compared with the `repositories.tags.name` attributes. If a match is
      not found and the registry is not registry.redhat.io, the check fails.
    - if the registry is registry.redhat.io, the check will skip the image.
- If the check fails use the point of failure to determine how to address the issue. 

For information on certifying images see: [Red Hat container certification](https://connect.redhat.com/partner-with-us/red-hat-container-certification)

### `chart-testing` v1.0

Chart testing runs the equivalant of `helm install ...` followed by `helm test...`. Try to run these independantly of 
the chart-verifier and make a note of any flags or overrides that must be set for them both to work. Ensure these 
values are set using chart-verifier flags when generating a report.

Also note that if chart-verifier flags are required for the chart-verifier chart-testing check to pass 
a verifier report must be included in the chart submission.

If the check fails due to a timeout, increase the timeout values. If increased timeouts are required for the chart-verifier chart-testing check to pass
a verifier report must be included in the chart submission. See: [Chart testing timeouts](./helm-chart-checks.md#chart-testing-timeouts)

Run the chart verifier and set log_ouput to true to get additional information:
```
$ podman run --rm -i \
          -e KUBECONFIG=/.kube/config \
          -v "${HOME}/.kube":/.kube:z \ 
          "quay.io/redhat-certification/chart-verifier" \
          verify -l \
          <chart-uri>
```

### `required-annotations-present` v1.0

Requires the following annotation to be present in chart.yaml:
- ```charts.openshift.io/name```

The value of the annotation will be used in the OpenShift catalogue as the name of the chart.


### `signature-is-valid` v1.0

This check requires that the public key provided to the chart verifier is from a user that has access to the signed chart. The check can fail for a variety of reasons, including:
- pgp public key file specified does not exist.
- pgp public key file is not an ascii public key file.
    - create using, for example: ```gpg --export -a <User-Name> > <public-key-file>```
      - User-Name is the user name of the secret key used to sign the chart.
- pgp public key file does not have access to the signed chart.
    - ensure the public key matches the secret key used to sign the chart. 
    

## Report related submission failures

### One or more mandatory checks have failed or are missing from the report.

Submission will fail if any [mandatory checks](./helm-chart-checks.md#default-set-of-checks-for-a-helm-chart) indicate failure or are absent from the report.

Regenerate the report running all tests and ensure they all pass.

If a check is failing and you are unsure as to why see [Trouble shooting check failures](#troubleshooting-check-failures)

### The digest in the report does not match the digest calculated for the submitted chart.

Common causes:

- The chart was updated after the report was generated.
    - Note: you can modify the webCatalogOnly attribute without affecting the report sha value.
- The Report was generated against a different form of the chart compared to the chart submitted.
    - For example report was generated from the chart source, but the chart tarball was used for submission.

For more information see [Verifier added annotations](./helm-chart-annotations.md#verifier-added-annotations)

### The certifiedOpenShiftVersions annotation does not contain a valid value.

This annotation must contain a current or recent OpenShift version. It is generally set by the chart-testing check
but this can fail if the role of the user who generated report does not have the required access.

For more information see [Verifier added annotations](./helm-chart-annotations.md#verifier-added-annotations)


## The chart uri is not a valid url.

For a report only submission the report must include a valid url for the chart.

For more information see [error-with-the-chart-url-when-submitting-report](https://github.com/openshift-helm-charts/charts/blob/main/docs/README.md#error-with-the-chart-url-when-submitting-report)

For more information see [Verifier added annotations](./helm-chart-annotations.md#verifier-added-annotations)
   
