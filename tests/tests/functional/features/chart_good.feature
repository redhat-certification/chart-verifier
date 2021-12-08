Feature: Chart  verification
    Partners or redhat or community can verify their charts by running the
    chart verifier against an error free chart.

    Examples:
        | image_type |
        | docker |
        | tarball |

    Scenario Outline: A chart provider verifies their chart using the chart verifier
        Given I would like to use the <type> profile
        Given I will provide a <location> of a <helm_chart>
        Given I will provide a <location> of an expected <report_info>
        When I run the <image_type> chart-verifier verify command against the chart to generate a report
        Then I should see the report-info from the generated report matching the expected report-info

        Examples:
            | type      | location                          | helm_chart              | report_info                     |
            | partner   | tests/charts/psql-service/0.1.8/  | src                     | partner-report-info.json   |
            | partner   | tests/charts/psql-service/0.1.9/  | psql-service-0.1.9.tgz  | partner-report-info.json   |
            | redhat    | tests/charts/psql-service/0.1.8/  | src                     | redhat-report-info.json    |
            | redhat    | tests/charts/psql-service/0.1.9/  | psql-service-0.1.9.tgz  | redhat-report-info.json    |
            | community | tests/charts/psql-service/0.1.8/  | src                     | community-report-info.json |
            | community | tests/charts/psql-service/0.1.9/  | psql-service-0.1.9.tgz  | community-report-info.json |


