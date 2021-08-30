## Creating a new chart-verifier release

Chart verifier - release creation is automated through a git hub workflow. To create a new release follow these steps:

1. Modify the [cmd/releases/release-info.json](https://github.com/redhat-certification/chart-verifier/blob/main/cmd/release/release_info.json) file with information about the new release, for example:
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

1. Create a PR which contains only [cmd/releases/release-info.json](https://github.com/redhat-certification/chart-verifier/blob/main/cmd/release/release_info.json)

1. The workflow will detect the file is changed and automatically merge the file and create the release if the following conditions are met:

   - The version specified is later that the version in the existing file.
   - The PR does not contain any other files.
   - The submitter has approval authority for the repository.
   - All tests pass. 

1. After merging the PR and creating the release the workflow continues in an attempt to add the ```latest``` tag to the new chart verifier image in quay. This may take a while.      

   - When the release is created, quay will detect the release and build a docker image for the release. 
        - see [chart verifier image tags in quay](https://quay.io/repository/redhat-certification/chart-verifier?tab=tags)
   - The workflow retries for up to 15 minutes for the docker image to appear in quay so that it can be linked. If this fails it can be done manually:
        1. Navigate to the [chart verifier image tags in quay](https://quay.io/repository/redhat-certification/chart-verifier?tab=tags)
            1. For the new image hit the options icon on the far right
            1. A drop down list appears, select "add a new tag" 
            1. A dialogue appears, enter the new tag name as "latest" 
            1. Quay detects latest is laready in use, select "Move" so it points to the new release.
   
Notes:
- To link the image to the ```latest``` tag in quay an auth token is required. This must be set as a repository secret "QUAY_AUTH_TOKEN"
    - For information on creating an auth token see: [Red Hat Quay API guide](https://access.redhat.com/documentation/en-us/red_hat_quay/3/html/red_hat_quay_api_guide/using_the_red_hat_quay_api) 
    - The workflow uses the [ChartVerfifierReleaser OAuth Applictaion](https://quay.io/organization/redhat-certification?tab=applications). 
      - If the application auth token is changed for any reason the repository secret must also be updated.
    