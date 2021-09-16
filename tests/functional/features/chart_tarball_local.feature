Feature: Chart local tarball verification
    Partners or redhat or community can verify their charts by running the
    chart verifier against a local chart tarball.

   Scenario: The partner IBM succesfully verifies their chart tarball using the chart verifier docker image
        Given: I am a Red Hat partner from IBM
            And: I would like to run chart verifier using the docker image
            And: I would like to use the partner profile
            And: I will provide a file location of an error free chart tarball
        When: I run the "chart-verifier verify" comand using the partner profile
        Then: I should see a report which passes all mandatory tests


    Scenario: RedHat succesfully verifies their chart tarball using the chart verifier docker image
        Given: I am a Red Hat chart provider
            And: I would like to run chart verifier using the docker image
            And: I would like to use the red hat profile
            And: I will provide a file location of an error free chart tarball
        When: I run the "chart-verifier verify" command using the red hat profile
        Then: I should see a report which passes all mandatory tests

    Scenario: Community partner succesfully verifies their chart tarball using the chart verifier docker image
       Given: I am a community chart provider
            And: I would like to run chart verifier using the docker image
            And: I would like to use the community profile
            And: I will provide a file location of an error free chart tarball
        When: I run the "chart-verifier verify" command using the community profile
        Then: I should see a report which passes all mandatory tests
