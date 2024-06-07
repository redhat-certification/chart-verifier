Feature: Chart verification
    Partners or redhat or community can verify their charts by running the
    chart verifier against an error free chart.

    @source
    Scenario Outline: A chart provider verifies their chart using the chart verifier
        Given I would like to use the <category> profile
        Given I will provide a <location> of a <helm_chart>
        Given I will provide a <location> of an expected <report_info>
        Given I will use the chart verifier <image_type> image
        Given The chart verifier version value
        When I run the chart-verifier verify command against the chart to generate a report
        Then I should see the report-info from the generated report matching the expected report-info
        Examples:
            | category  | location                          | helm_chart              | report_info                | image_type |
            # test with tarball
            | partner   | tests/charts/psql-service/0.1.8/  | src                     | partner-report-info.json   | tarball    | 
            | partner   | tests/charts/psql-service/0.1.9/  | psql-service-0.1.9.tgz  | partner-report-info.json   | tarball    |
            | redhat    | tests/charts/psql-service/0.1.8/  | src                     | redhat-report-info.json    | tarball    |
            | redhat    | tests/charts/psql-service/0.1.9/  | psql-service-0.1.9.tgz  | redhat-report-info.json    | tarball    |
            | community | tests/charts/psql-service/0.1.8/  | src                     | community-report-info.json | tarball    |
            | community | tests/charts/psql-service/0.1.9/  | psql-service-0.1.9.tgz  | community-report-info.json | tarball    |
            # Testing with Podman
            | partner   | tests/charts/psql-service/0.1.8/  | src                     | partner-report-info.json   | podman     |
            | partner   | tests/charts/psql-service/0.1.9/  | psql-service-0.1.9.tgz  | partner-report-info.json   | podman     |
            | redhat    | tests/charts/psql-service/0.1.8/  | src                     | redhat-report-info.json    | podman     |
            | redhat    | tests/charts/psql-service/0.1.9/  | psql-service-0.1.9.tgz  | redhat-report-info.json    | podman     |
            | community | tests/charts/psql-service/0.1.8/  | src                     | community-report-info.json | podman     |
            | community | tests/charts/psql-service/0.1.9/  | psql-service-0.1.9.tgz  | community-report-info.json | podman     |

    @signed
    Scenario Outline: A chart provider verifies their signed chart using the chart verifier
        Given I would like to use the <category> profile
        Given I will provide a <location> of a <helm_chart>
        Given I will provide a <location> of an expected <report_info>
        Given I will use the chart verifier <image_type> image
        Given I will provide a <location> of a <public_key> to verify the signature
        Given The chart verifier version value
        When I run the chart-verifier verify command against the signed chart to generate a report
        Then I should see the report-info from the report for the signed chart matching the expected report-info
        Examples:
            | category  | location                           | helm_chart               | report_info                | public_key                  | image_type |
            # test with tarball
            | partner   | tests/charts/psql-service/0.1.11/  | psql-service-0.1.11.tgz  | partner-report-info.json   | psql-service-0.1.11.tgz.key | tarball    |
            | redhat    | tests/charts/psql-service/0.1.11/  | psql-service-0.1.11.tgz  | redhat-report-info.json    | psql-service-0.1.11.tgz.key | tarball    |
            # test with podman
            | partner   | tests/charts/psql-service/0.1.11/  | psql-service-0.1.11.tgz  | partner-report-info.json   | psql-service-0.1.11.tgz.key | podman     |
            | redhat    | tests/charts/psql-service/0.1.11/  | psql-service-0.1.11.tgz  | redhat-report-info.json    | psql-service-0.1.11.tgz.key | podman     |
