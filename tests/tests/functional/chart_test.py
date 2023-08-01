"""This file contains the tests for chart-verifier. It follows behavioral driven
development using pytest-bdd.

Following environment variables are expected to be set in order to run the tests:
- VERIFIER_TARBALL_NAME: Path to the tarball to use for the tarball tests.
- PODMAN_IMAGE_TAG: Tag of the container image to use to run the podman tests.
- KUBECONFIG: Path to the kubeconfig file to use to connect to the OCP cluster to
              deploy the charts on.
"""


from pytest_bdd import scenario, given, when, then, parsers
import os
import subprocess
import pytest
import json
import tarfile
from deepdiff import DeepDiff

import yaml

try:
    from yaml import CLoader as Loader
except ImportError:
    from yaml import Loader

REPORT_ANNOTATIONS = "annotations"
REPORT_RESULTS = "results"
REPORT_DIGESTS = "digests"
REPORT_METADATA = "metadata"


@scenario(
    "features/chart_good.feature",
    "A chart provider verifies their chart using the chart verifier",
)
def test_chart_source():
    pass


@scenario(
    "features/chart_good.feature",
    "A chart provider verifies their signed chart using the chart verifier",
)
def test_chart_signed():
    pass


@given(
    parsers.parse("I would like to use the <type> profile"),
    target_fixture="profile_type",
)
def profile_type(type):
    return type


@given(
    parsers.parse("I will provide a <location> of a <helm_chart>"),
    target_fixture="chart_location",
)
def chart_location(location, helm_chart):
    return os.path.join(location, helm_chart)


@given(
    parsers.parse("I will provide a <location> of an expected <report_info>"),
    target_fixture="report_info_location",
)
def report_info_location(location, report_info):
    return os.path.join(location, report_info)


@given(
    parsers.parse(
        "I will provide a <location> of a <public_key> to verify the signature"
    ),
    target_fixture="public_key_location",
)
def public_key_location(location, public_key):
    return os.path.join(location, public_key)


@given(
    parsers.parse("I will use the chart verifier <image_type> image"),
    target_fixture="image_type",
)
def image_type(image_type):
    return image_type


@given("The chart verifier version value", target_fixture="verifier_version")
def verifier_version(image_type):
    """Get the version of the chart verifier tool used to produce and verify reports.

    This output comes directly from the output of `chart-verifier version`, which
    is then normalized to match what we would expect to find in a report.

    Args:
        image_type (string): How chart verifier will run. Options: tarball, podman

    Returns:
        str: a normalized semantic version, like 0.0.0
    """
    if image_type == "tarball":
        tarball_name = os.environ.get("VERIFIER_TARBALL_NAME")
        print(f"\nRun version using tarbal {tarball_name}")
        return run_version_tarball_image(tarball_name)

    # Podman
    image_tag = os.environ.get("PODMAN_IMAGE_TAG")
    if not image_tag:
        image_tag = "main"
    image_name = "quay.io/redhat-certification/chart-verifier"
    print(f"\nRun version using podman image {image_name}:{image_tag}")
    return run_version_podman_image(image_name, image_tag)


@when(
    parsers.parse(
        "I run the chart-verifier verify command against the chart to generate a report"
    ),
    target_fixture="run_verify",
)
def run_verify(image_type, profile_type, chart_location):
    print(
        f"\nrun {image_type} verifier verify  with profile : {profile_type},"
        f" and chart: {chart_location}"
    )
    return run_verifier(image_type, profile_type, chart_location, "verify")


@when(
    parsers.parse(
        "I run the chart-verifier verify command against the signed chart to generate "
        "a report"
    ),
    target_fixture="run_signed_verify",
)
def run_signed_verify(image_type, profile_type, chart_location, public_key_location):
    print(
        f"\nrun {image_type} verifier verify  with profile : {profile_type}, "
        f"and signed chart: {chart_location}"
    )
    return run_verifier(
        image_type, profile_type, chart_location, "verify", public_key_location
    )


def run_report(image_type, profile_type, report_location):
    print(
        f"\nrun {image_type} verifier report  with profile : {profile_type}, "
        f"and chart: {report_location}"
    )
    return run_verifier(image_type, profile_type, report_location, "report")


def run_verifier(
    image_type, profile_type, target_location, command, pgp_key_location=None
):
    """Run chart-verifier command and capture output

    Args:
        image_type (str): How to run chart-verifier. Options: tarball, podman.
        profile_type (str): Profile to use. Options: partner, redhat, community.
        target_location (str): Path to the Helm chart to run verify against.
        command (str): Command to run, either verify or report.
        pgp_key_location (str, optional): Path to the GPG public key of the key used to
                                          sign the chart.

    Returns:
        str: The output provided by the command.
    """
    if image_type == "tarball":
        tarball_name = os.environ.get("VERIFIER_TARBALL_NAME")
        print(f"\nRun {command} using tarball: {tarball_name}")
        if command == "verify":
            return run_verify_tarball_image(
                tarball_name, profile_type, target_location, pgp_key_location
            )
        else:
            return run_report_tarball_image(tarball_name, profile_type, target_location)

    # Podman
    image_tag = os.environ.get("PODMAN_IMAGE_TAG")
    if not image_tag:
        image_tag = "main"
    image_name = "quay.io/redhat-certification/chart-verifier"
    print(f"\nRun {command} using podman image {image_name}:{image_tag}")
    if command == "verify":
        return run_verify_podman_image(
            image_name, image_tag, profile_type, target_location, pgp_key_location
        )
    else:
        return run_report_podman_image(
            image_name, image_tag, profile_type, target_location
        )


def run_version_tarball_image(tarball_name):
    tar = tarfile.open(tarball_name, "r:gz")
    tar.extractall(path="./test_verifier")
    out = subprocess.run(
        ["./test_verifier/chart-verifier", "version", "--as-data"], capture_output=True
    )
    return normalize_version(out.stdout.decode("utf-8"))


def normalize_version(version):
    """Extract normalized version from JSON data

    Args:
        version (str): the version and commit ID in JSON

    Returns:
        str: a normalized semver like 0.0.0.
    """
    print(f"version input to normalize_version function is: {version}")
    version_dict = json.loads(version)
    return version_dict["version"]


def run_version_podman_image(verifier_image_name, verifier_image_tag):
    """Run chart verifier's version command in Podman."""
    out = subprocess.run(
        [
            "podman",
            "run",
            "--rm",
            f"{verifier_image_name}:{verifier_image_tag}",
            "version",
            "--as-data",
        ],
        capture_output=True,
    )
    return normalize_version(out.stdout.decode("utf-8"))


def run_verify_tarball_image(
    tarball_name, profile_type, chart_location, pgp_key_location=None
):
    print(f"Run tarball image from {tarball_name}")

    tar = tarfile.open(tarball_name, "r:gz")

    tar.extractall(path="./test_verifier")

    if pgp_key_location:
        out = subprocess.run(
            [
                "./test_verifier/chart-verifier",
                "verify",
                "--set",
                f"profile.vendorType={profile_type}",
                "--pgp-public-key",
                pgp_key_location,
                chart_location,
            ],
            capture_output=True,
        )
    else:
        out = subprocess.run(
            [
                "./test_verifier/chart-verifier",
                "verify",
                "--set",
                f"profile.vendorType={profile_type}",
                chart_location,
            ],
            capture_output=True,
        )

    return out.stdout.decode("utf-8")


def run_report_tarball_image(tarball_name, profile_type, chart_location):
    print(f"Run tarball image from {tarball_name}")

    tar = tarfile.open(tarball_name, "r:gz")

    tar.extractall(path="./test_verifier")

    out = subprocess.run(
        [
            "./test_verifier/chart-verifier",
            "report",
            "all",
            "--set",
            f"profile.vendorType={profile_type}",
            chart_location,
        ],
        capture_output=True,
    )

    return out.stdout.decode("utf-8")


def run_verify_podman_image(
    verifier_image_name,
    verifier_image_tag,
    profile_type,
    chart_location,
    pgp_key_location=None,
):
    print(f"Run podman image - {verifier_image_name}:{verifier_image_tag}")
    kubeconfig = os.environ.get("KUBECONFIG")
    if not kubeconfig:
        return "FAIL: missing KUBECONFIG environment variable"

    if chart_location.startswith("http:/") or chart_location.startswith("https:/"):
        if pgp_key_location:
            out = subprocess.run(
                [
                    "podman",
                    "run",
                    "-v",
                    f"{kubeconfig}:/kubeconfig:z",
                    "-e",
                    "KUBECONFIG=/kubeconfig",
                    "--rm",
                    f"{verifier_image_name}:{verifier_image_tag}",
                    "verify",
                    "--set",
                    f"profile.vendortype={profile_type}",
                    "--pgp-public-key",
                    public_key_location,
                    chart_location,
                ],
                capture_output=True,
            )
        else:
            out = subprocess.run(
                [
                    "podman",
                    "run",
                    "-v",
                    f"{kubeconfig}:/kubeconfig:z",
                    "-e",
                    "KUBECONFIG=/kubeconfig",
                    "--rm",
                    f"{verifier_image_name}:{verifier_image_tag}",
                    "verify",
                    "--set",
                    f"profile.vendortype={profile_type}",
                    chart_location,
                ],
                capture_output=True,
            )
    else:
        chart_directory = os.path.dirname(os.path.abspath(chart_location))
        chart_name = os.path.basename(os.path.abspath(chart_location))
        if pgp_key_location:
            pgp_key_name = os.path.basename(os.path.abspath(pgp_key_location))
            out = subprocess.run(
                [
                    "podman",
                    "run",
                    "-v",
                    f"{chart_directory}:/charts:z",
                    "-v",
                    f"{kubeconfig}:/kubeconfig:z",
                    "-e",
                    "KUBECONFIG=/kubeconfig",
                    "--rm",
                    f"{verifier_image_name}:{verifier_image_tag}",
                    "verify",
                    "--set",
                    f"profile.vendortype={profile_type}",
                    "--pgp-public-key",
                    f"/charts/{pgp_key_name}",
                    f"/charts/{chart_name}",
                ],
                capture_output=True,
            )
        else:
            out = subprocess.run(
                [
                    "podman",
                    "run",
                    "-v",
                    f"{chart_directory}:/charts:z",
                    "-v",
                    f"{kubeconfig}:/kubeconfig:z",
                    "-e",
                    "KUBECONFIG=/kubeconfig",
                    "--rm",
                    f"{verifier_image_name}:{verifier_image_tag}",
                    "verify",
                    "--set",
                    f"profile.vendortype={profile_type}",
                    f"/charts/{chart_name}",
                ],
                capture_output=True,
            )

    return out.stdout.decode("utf-8")


def run_report_podman_image(
    verifier_image_name, verifier_image_tag, profile_type, report_location
):
    print(f"Run podman image - {verifier_image_name}:{verifier_image_tag}")

    report_directory = os.path.dirname(os.path.abspath(report_location))
    report_name = os.path.basename(os.path.abspath(report_location))
    out = subprocess.run(
        [
            "podman",
            "run",
            "-v",
            f"{report_directory}:/reports:z",
            "--rm",
            f"{verifier_image_name}:{verifier_image_tag}",
            "report",
            "all",
            "--set",
            f"profile.vendortype={profile_type}",
            f"/reports/{report_name}",
        ],
        capture_output=True,
    )

    return out.stdout.decode("utf-8")


@then(
    "I should see the report-info from the report for the signed chart matching the "
    "expected report-info"
)
def signed_chart_report(
    run_signed_verify, profile_type, report_info_location, image_type, verifier_version
):
    check_report(
        run_signed_verify,
        profile_type,
        report_info_location,
        image_type,
        verifier_version,
    )


@then(
    "I should see the report-info from the generated report matching the expected "
    "report-info"
)
def chart_report(
    run_verify, profile_type, report_info_location, image_type, verifier_version
):
    check_report(
        run_verify, profile_type, report_info_location, image_type, verifier_version
    )


def check_report(
    verify_result, profile_type, report_info_location, image_type, verifier_version
):
    """Compares the output of chart-verifier against a pre-generated report, containing
    the expected results.

    Fails the tests if a difference is found.

    Args:
        verify_result (str): Output of chart-verifier command.
        profile_type (str): Profile to use. Options: partner, redhat, community.
        report_info_location (str): Path to the pre-generated report.
        image_type (str): How to run chart-verifier. Options: tarball, podman.
        verifier_version (str): Normalized version as given by chart-verifier.
    """
    if verify_result.startswith("FAIL"):
        pytest.fail(f"FAIL some tests failed: {verify_result}")

    print(f"Report data:\n{verify_result}\ndone")

    report_data = yaml.load(verify_result, Loader=Loader)

    test_passed = True

    report_verifier_version = report_data["metadata"]["tool"]["verifier-version"]
    if report_verifier_version != verifier_version:
        print(
            "FAIL: verifier-version found in report does not match tool version. "
            f"Expected {verifier_version}, but report has {report_verifier_version}"
        )
        test_passed = False

    report_vendor_type = report_data["metadata"]["tool"]["profile"]["VendorType"]
    if report_vendor_type != profile_type:
        print(
            "FAIL: profiles do not match. "
            f"Expected {profile_type}, but report has {report_vendor_type}"
        )
        test_passed = False

    chart_name = report_data["metadata"]["chart"]["name"]
    chart_version = report_data["metadata"]["chart"]["version"]

    report_name = f"{profile_type}-{chart_name}-{chart_version}-report.yaml"
    report_dir = os.path.join(os.getcwd(), "test-reports")
    report_path = os.path.join(report_dir, report_name)

    if not os.path.isdir(report_dir):
        os.makedirs(report_dir)

    print(f"Report path : {report_path}")

    with open(report_path, "w") as fd:
        fd.write(verify_result)

    expected_reports_file = open(
        report_info_location,
    )
    expected_reports = json.load(expected_reports_file)

    test_report_data = run_report(image_type, profile_type, report_path)
    print(f"test_report_data : \n{test_report_data}")
    test_report = json.loads(test_report_data)

    results_diff = DeepDiff(
        expected_reports[REPORT_RESULTS], test_report[REPORT_RESULTS], ignore_order=True
    )
    if results_diff:
        print(f"difference found in results : {results_diff}")
        test_passed = False

    expected_annotations = {}
    for report_annotation in expected_reports[REPORT_ANNOTATIONS]:
        expected_annotations[report_annotation["name"]] = report_annotation["value"]

    tested_annotations = {}
    for report_annotation in test_report[REPORT_ANNOTATIONS]:
        tested_annotations[report_annotation["name"]] = report_annotation["value"]

    missing_annotations = set(expected_annotations.keys()) - set(
        tested_annotations.keys()
    )
    if missing_annotations:
        test_passed = False
        print(f"Missing annotations: {missing_annotations}")

    extra_annotations = set(tested_annotations.keys()) - set(
        expected_annotations.keys()
    )
    if extra_annotations:
        test_passed = False
        print(f"extra annotations: {extra_annotations}")

    differing_annotations = {
        "charts.openshift.io/lastCertifiedTimestamp",
        "charts.openshift.io/testedOpenShiftVersion",
    }

    for annotation in expected_annotations.keys():
        if annotation not in differing_annotations:
            if expected_annotations[annotation] != tested_annotations[annotation]:
                test_passed = False
                print(
                    f"{annotation} has different content, "
                    f"expected: {expected_annotations[annotation]}, "
                    f"got: {tested_annotations[annotation]}"
                )

    digests_diff = DeepDiff(
        expected_reports[REPORT_DIGESTS], test_report[REPORT_DIGESTS], ignore_order=True
    )
    if digests_diff:
        print(f"difference found in digests : {digests_diff}")
        test_passed = False

    differing_metadata = {"chart-uri"}
    expected_metadata = expected_reports[REPORT_METADATA]
    test_metadata = test_report[REPORT_METADATA]
    for metadata in expected_metadata.keys():
        if metadata not in differing_metadata:
            metadata_diff = DeepDiff(
                expected_metadata[metadata], test_metadata[metadata], ignore_order=True
            )
            if metadata_diff:
                test_passed = False
                print(f"difference found in {metadata} metadata : {metadata_diff}")

    if not test_passed:
        pytest.fail("FAIL differences found in reports")
