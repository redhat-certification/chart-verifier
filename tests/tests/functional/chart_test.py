
from pytest_bdd import scenario, given, when, then, parsers
import os
import subprocess
import sys
import docker
import pytest
import json
import tarfile
from deepdiff import DeepDiff

sys.path.append('./scripts/src/')
from report import report_info

import yaml
try:
    from yaml import CLoader as Loader, CDumper as Dumper
except ImportError:
    from yaml import Loader, Dumper


@scenario(
    "features/chart_good.feature",
    "A chart provider verifies their chart using the chart verifier"
)
def test_chart_source():
    pass

@given(parsers.parse("I would like to use the <type> profile"),target_fixture="profile_type")
def profile_type(type):
    return type

@given(parsers.parse("I will provide a <location> of a <helm_chart>"),target_fixture="chart_location")
def chart_location(location,helm_chart):
    return os.path.join(location,helm_chart)

@given(parsers.parse("I will provide a <location> of an expected <report_info>"),target_fixture="report_info_location")
def report_info_location(location,report_info):
    return os.path.join(location,report_info)

@when(parsers.parse("I run the <image_type> chart-verifier verify command against the chart to generate a report"),target_fixture="run_verifier")
def run_verifier(image_type, profile_type, chart_location):
    print(f"\nrun {image_type} verifier  with profile : {profile_type}, and chart: {chart_location}")

    if image_type == "tarball":
        tarball_name = os.environ.get("VERIFIER_TARBALL_NAME")
        return run_tarball_image(tarball_name,profile_type,chart_location)
    elif image_type == "docker":
        image_tag  =  os.environ.get("VERIFER_IMAGE_TAG")
        if not image_tag:
            image_tag = "main"
        image_name =  "quay.io/redhat-certification/chart-verifier"
        return run_docker_image(image_name,image_tag,profile_type,chart_location)
    else:
        image_tag = os.environ.get("PODMAN_IMAGE_TAG")
        if not image_tag:
            image_tag = "main"
        image_name =  "quay.io/redhat-certification/chart-verifier"
        return run_podman_image(image_name,image_tag,profile_type,chart_location)

def run_docker_image(verifier_image_name,verifier_image_tag,profile_type, chart_location):

    client = docker.from_env()

    try:
        verifier_image=client.images.pull(verifier_image_name,tag=verifier_image_tag)
    except docker.errors.APIError as exc:
        print(f'Error from docker loading image: {verifier_image_name}:{verifier_image_tag}')
        return f"FAIL pulling image : docker.errors.APIError: {exc.args}"

    os.environ["VERIFIER_IMAGE"] = f"{verifier_image_name}:{verifier_image_tag}"

    docker_command = "verify "
    local_chart = False
    if chart_location.startswith('http:/') or chart_location.startswith('https:/'):
        docker_command = docker_command + chart_location
    else:
        if os.path.exists(chart_location):
            docker_command = docker_command + f"/charts/{os.path.basename(chart_location)}"
            local_chart = True
        else:
            return f"FAIL: chart does not exist: {os.path.abspath(chart_location)}"

    if profile_type:
        docker_command = docker_command + " --set profile.vendorType=" + profile_type

    print(f'docker command: {docker_command}')

    kubeconfig = os.environ.get("KUBECONFIG")
    if not kubeconfig:
        return "FAIL: missing KUBECONFIG environment variable"

    docker_volumes = {kubeconfig: {'bind': '/kubeconfig', 'mode': 'ro'}}
    docker_environment = {"KUBECONFIG": '/kubeconfig'}

    try:
        if local_chart:
            chart_directory = os.path.dirname(os.path.abspath(chart_location))
            docker_volumes[chart_directory] = {'bind': '/charts/', 'mode': 'rw'}

        output = client.containers.run(verifier_image,docker_command,stdin_open=True,tty=True,stdout=True,volumes=docker_volumes,environment=docker_environment)

    except docker.errors.ContainerError as exc:
        return f"FAIL: docker.errors.ContainerError: {exc.args}"
    except docker.errors.ImageNotFound as exc:
        return f"FAIL: docker.errors.ImageNotFound: {exc.args}"
    except docker.errors.APIError as exc:
        return f"FAIL: docker.errors.APIError: {exc.args}"

    if not output:
        return f"FAIL: no report produced : {docker_command}"

    return output.decode("utf-8")

def run_tarball_image(tarball_name,profile_type, chart_location):
    print(f"Run tarball image from {tarball_name}")

    tar = tarfile.open(tarball_name, "r:gz")

    tar.extractall(path="./test_verifier")

    out = subprocess.run(["./test_verifier/chart-verifier","verify","--set",f"profile.vendorType={profile_type}",chart_location],capture_output=True)

    return out.stdout.decode("utf-8")

def run_podman_image(verifier_image_name,verifier_image_tag,profile_type, chart_location):

    print(f"Run podman image - {verifier_image_name}:{verifier_image_tag}")
    kubeconfig = os.environ.get("KUBECONFIG")
    if not kubeconfig:
        return "FAIL: missing KUBECONFIG environment variable"

    if chart_location.startswith('http:/') or chart_location.startswith('https:/'):
        out = subprocess.run(["podman", "run", "-v", f"{kubeconfig}:/kubeconfig", "-e", "KUBECONFIG=/kubeconfig", "--rm",
                          f"{verifier_image_name}:{verifier_image_tag}", "verify", "--set", f"profile.vendortype={profile_type}", chart_location], capture_output=True)
    else:
        chart_directory = os.path.dirname(os.path.abspath(chart_location))
        chart_name = os.path.basename(os.path.abspath(chart_location))
        out = subprocess.run(["podman", "run", "-v", f"{chart_directory}:/charts:z", "-v", f"{kubeconfig}:/kubeconfig", "-e", "KUBECONFIG=/kubeconfig", "--rm",
                              f"{verifier_image_name}:{verifier_image_tag}", "verify", "--set", f"profile.vendortype={profile_type}", f"/charts/{chart_name}"], capture_output=True)

    return out.stdout.decode("utf-8")

@then("I should see the report-info from the generated report matching the expected report-info")
def check_report(run_verifier, profile_type, report_info_location):

    if run_verifier.startswith("FAIL"):
        pytest.fail(f'FAIL some tests failed: {run_verifier}')

    print(f"Report data:\n{run_verifier}\ndone")

    report_data = yaml.load(run_verifier, Loader=Loader)

    test_passed = True

    report_vendor_type = report_data["metadata"]["tool"]["profile"]["VendorType"]
    if report_vendor_type != profile_type:
        print(f"FAIL: profiles do not match. Expected {profile_type}, but report has {report_vendor_type}")
        test_passed = False

    report_version = report_data["metadata"]["tool"]["profile"]["version"]

    chart_name =  report_data["metadata"]["chart"]["name"]
    chart_version = report_data["metadata"]["chart"]["version"]

    report_name = f'{profile_type}-{chart_name}-{chart_version}-report.yaml'
    report_dir = os.path.join(os.getcwd(),"test-reports")
    report_path = os.path.join(report_dir,report_name)

    if not os.path.isdir(report_dir):
        os.makedirs(report_dir)

    print(f'Report path : {report_path}')

    with open(report_path, "w") as fd:
        fd.write(run_verifier)

    test_reports = report_info.get_all_reports(report_path,report_vendor_type,report_version)

    expected_reports_file = open(report_info_location,)
    expected_reports = json.load(expected_reports_file)

    results_diff = DeepDiff(expected_reports[report_info.REPORT_RESULTS],test_reports[report_info.REPORT_RESULTS],ignore_order=True)
    if results_diff:
        print(f"difference found in results : {results_diff}")
        test_passed = False

    expected_annotations = {}
    for report_annotation in expected_reports[report_info.REPORT_ANNOTATIONS]:
        expected_annotations[report_annotation["name"]] = report_annotation["value"]

    tested_annotations = {}
    for report_annotation in test_reports[report_info.REPORT_ANNOTATIONS]:
        tested_annotations[report_annotation["name"]] = report_annotation["value"]

    missing_annotations = set(expected_annotations.keys()) - set(tested_annotations.keys())
    if missing_annotations:
        test_passed = False
        print(f"Missing annotations: {missing_annotations}")

    extra_annotations = set(tested_annotations.keys()) - set(expected_annotations.keys())
    if extra_annotations:
        test_passed = False
        print(f"extra annotations: {extra_annotations}")

    differing_annotations = {"charts.openshift.io/lastCertifiedTimestamp",
                             "charts.openshift.io/testedOpenShiftVersion"}

    for annotation in expected_annotations.keys():
        if not annotation in differing_annotations:
            if expected_annotations[annotation] != tested_annotations[annotation]:
                test_passed = False
                print(f"{annotation} has different content, expected: {expected_annotations[annotation]}, got: {tested_annotations[annotation]}")

    digests_diff = DeepDiff(expected_reports[report_info.REPORT_DIGESTS],test_reports[report_info.REPORT_DIGESTS],ignore_order=True)
    if digests_diff:
        print(f"difference found in digests : {digests_diff}")
        test_passed = False

    differing_metadata = {"chart-uri"}
    expected_metadata = expected_reports[report_info.REPORT_METADATA]
    test_metadata = test_reports[report_info.REPORT_METADATA]
    for metadata in expected_metadata.keys():
        if not metadata in differing_metadata:
            metadata_diff = DeepDiff(expected_metadata[metadata],test_metadata[metadata],ignore_order=True)
            if metadata_diff:
                test_passed = False
                print(f"difference found in {metadata} metadata : {metadata_diff}")

    if not test_passed:
        pytest.fail('FAIL differences found in reports')


