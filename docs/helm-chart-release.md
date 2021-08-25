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
    
1. Once the release is created, quay will detect the release and build a docker image for the release. 
    - see [chart verifier tags in quay](https://quay.io/repository/redhat-certification/chart-verifier?tab=tags)

1. Once the image is in quay the latest image mut be reset to point to the new release:
   1. For the new image hit the options icon on the far right
   1. A drop down list appears, select "add a new tag" 
   1. A dialogue appears, enter the new tag name as "latest" 
   1. Quay detects latest is laready in use, select "Move" so it points to the new release.
   
    Note: it is intended to automate this step in the future. 