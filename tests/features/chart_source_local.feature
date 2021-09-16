Feature: Chart local source verification
    Partners or redhat or community can verify their charts by running the
    chart verifier against a local chart source

    Scenario: The partner IBM succesfully verifies their chart using the chart verifier docker image
        Given IBM runs that chart verifier using the docker image and using the partner profile
        And IBM has specfied an error free local chart source
        IBM sees a report which passes all mandatory tests

    Scenario: RedHat succesfully verifies their chart using the chart verifier docker image
        Given RedHat runs that chart verifier using the docker image and using the redhat profile
        And RedHat has specfied an error free local chart source
        RedHat sees a report which passes all mandatory tests

    Scenario: Community partner succesfully verifies their chart using the chart verifier docker image
        Given Community partner runs that chart verifier using the docker image and using the community profile
        And Community partner has specfied a local chart source
        Community partner sees a report which passes all mandatory tests
