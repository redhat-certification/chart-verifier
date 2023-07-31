## Creating a new chart-verifier release

Chart verifier - release creation is automated through a GitHub workflow. To create a new release follow these steps:

1. Modify the [pkg/chartverifier/version/version_info.json](https://github.com/redhat-certification/chart-verifier/blob/main/pkg/chartverifier/version/version_info.json) file with information about the new release, for example:
   ```
   {
       "version": "1.2.0",
       "quay-image":  "quay.io/redhat-certification/chart-verifier",
       "release-info": [
           "Additions to report command metadata output: #174",
           "New report command: #170"
       ]
   }
   ```
    - version : set the new chart-verifer version.
    - quay-image : set to the name of the image (only used for testing).
    - release-info : list of significant PR's in the release

1. Create a PR which contains only [pkg/chartverifier/version/version_info.json](https://github.com/redhat-certification/chart-verifier/blob/main/pkg/chartverifier/version/version_info.json)

1. The workflow will detect the file is changed and automatically merge the file and create the release if the following conditions are met:

   - The version specified is later that the version in the existing file.
   - The PR does not contain any other files.
   - The submitter has approval authority for the repository.
   - All tests pass. 

1. After the PR is merged, the "release.yaml" workflow runs and create the GitHub release along with all its assets:

    - A container image is built and pushed to Quay under two tags: ```latest``` and one corresponding to the version number (```x.y.z``` without leading ```v``).
    - A tarball is created and attached to the GitHub release.

Notes:
- To push the images to Quay, an auth token is required. This must be set as a repository secret "QUAY_AUTH_TOKEN"
    - For information on creating an auth token see: [Red Hat Quay API guide](https://access.redhat.com/documentation/en-us/red_hat_quay/3/html/red_hat_quay_api_guide/using_the_red_hat_quay_api) 
    - The workflow uses the [ChartVerfifierReleaser OAuth Applictaion](https://quay.io/organization/redhat-certification?tab=applications). 
      - If the application auth token is changed for any reason the repository secret must also be updated.
