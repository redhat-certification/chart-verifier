
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

@given(parsers.parse("I would like to use the <type> profile"),target_fixture="profile_type")
def profile_type(type):
    return type

@given(parsers.parse("I will provide a <location> of a <helm_chart>"),target_fixture="chart_location")
def chart_location(location,helm_chart):
    return os.path.join(location,helm_chart)

@given(parsers.parse("I will provide a <location> of an expected <report_info>"),target_fixture="report_info_location")
def report_info_location(location,report_info):
    return os.path.join(location,report_info)

@given(parsers.parse("I will provide a <location> of a <public_key> to verify the signature"),target_fixture="public_key_location")
def public_key_location(location,public_key):
    return os.path.join(location,public_key)

@given(parsers.parse("I will use the chart verifier <image_type> image"),target_fixture="image_type")
def image_type(image_type):
    return image_type

@given('The chart verifier version value',target_fixture='verifier_version')
def verifier_version(image_type):
    """Get the version of the chart verifier tool used to produce and verify reports.

    This output comes directly from the output of `chart-verifier version`, which
    is the normalized to match what we would expect to find in a report.

    Parameters:
    image_type (string): How chart verifier will run. Options: tarball, podman, docker

    Returns:
    string: a normalized semantic version, like 0.0.0
    """
    if image_type == "tarball":
        tarball_name = os.environ.get("VERIFIER_TARBALL_NAME")
        print(f"\nRun version using tarbal {tarball_name}")
        return run_version_tarball_image(tarball_name)
    elif image_type == "podman":
        image_tag = os.environ.get("PODMAN_IMAGE_TAG")
        if not image_tag:
            image_tag = "main"
        image_name =  "quay.io/redhat-certification/chart-verifier"
        print(f"\nRun version using podman image {image_name}:{image_tag}")
        return run_version_podman_image(image_name,image_tag)
    else: # Fallback to Docker.
        image_tag  =  os.environ.get("VERIFIER_IMAGE_TAG")
        if not image_tag:
            image_tag = "main"
        image_name =  "quay.io/redhat-certification/chart-verifier"
        print(f"\nRun version using docker image {image_name}:{image_tag}")
        return run_version_docker_image(image_name, image_tag)

@when(parsers.parse("I run the chart-verifier verify command against the chart to generate a report"),target_fixture="run_verify")
def run_verify(image_type, profile_type, chart_location):
    print(f"\nrun {image_type} verifier verify  with profile : {profile_type}, and chart: {chart_location}")
    return run_verifier(image_type, profile_type, chart_location,"verify")

@when(parsers.parse("I run the chart-verifier verify command against the signed chart to generate a report"),target_fixture="run_signed_verify")
def run_signed_verify(image_type, profile_type, chart_location, public_key_location):
    print(f"\nrun {image_type} verifier verify  with profile : {profile_type}, and signed chart: {chart_location}")
    return run_verifier(image_type, profile_type, chart_location,"verify",public_key_location)

def run_report(image_type, profile_type, report_location):
    print(f"\nrun {image_type} verifier report  with profile : {profile_type}, and chart: {report_location}")
    return run_verifier(image_type, profile_type, report_location,"report")

def run_verifier(image_type, profile_type, target_location, command,pgp_key_location=None):

    if image_type == "tarball":
        tarball_name = os.environ.get("VERIFIER_TARBALL_NAME")
        print(f"\nRun {command} using tarball: {tarball_name}")
        if command == "verify":
            return run_verify_tarball_image(tarball_name,profile_type,target_location,pgp_key_location)
        else:
            return run_report_tarball_image(tarball_name,profile_type,target_location)
    elif image_type == "podman":
        image_tag = os.environ.get("PODMAN_IMAGE_TAG")
        if not image_tag:
            image_tag = "main"
        image_name =  "quay.io/redhat-certification/chart-verifier"
        print(f"\nRun {command} using podman image {image_name}:{image_tag}")
        if command == "verify":
            return run_verify_podman_image(image_name,image_tag,profile_type,target_location,pgp_key_location)
        else:
            return run_report_podman_image(image_name,image_tag,profile_type,target_location)
    else:
        image_tag  =  os.environ.get("VERIFIER_IMAGE_TAG")
        if not image_tag:
            image_tag = "main"
        image_name =  "quay.io/redhat-certification/chart-verifier"
        print(f"\nRun {command} using docker image {image_name}:{image_tag}")
        if command == "verify":
            return run_verify_docker_image(image_name,image_tag,profile_type,target_location,pgp_key_location)
        else:
            return run_report_docker_image(image_name,image_tag,profile_type,target_location)

def run_verify_docker_image(verifier_image_name,verifier_image_tag,profile_type, chart_location, pgp_key_location=None):

    client = docker.from_env()

    try:
        verifier_image=client.images.get(f'{verifier_image_name}:{verifier_image_tag}')
    except docker.errors.APIError as exc:
        try:
            print(f'Error from docker get image: {verifier_image_name}:{verifier_image_tag}')
            verifier_image=client.images.pull(verifier_image_name,tag=verifier_image_tag)
        except:
            print(f'Error from docker pulling image: {verifier_image_name}:{verifier_image_tag}')
            return f"FAIL getting image : docker.errors.APIError: {exc.args}"

    docker_command = "verify"

    if pgp_key_location:
        if pgp_key_location.startswith('http:/') or pgp_key_location.startswith('https:/'):
            docker_command = f"{docker_command}  --pgp-public-key {pgp_key_location}"
        else:
            if os.path.exists(pgp_key_location):
                docker_command = f"{docker_command} --pgp-public-key /charts/{os.path.basename(pgp_key_location)}"
            else:
                return f"FAIL: pgp public key not exist: {os.path.abspath(pgp_key_location)}"

    local_chart = False
    if chart_location.startswith('http:/') or chart_location.startswith('https:/'):
        docker_command = f"{docker_command}  {chart_location}"
    else:
        if os.path.exists(chart_location):
            docker_command = f"{docker_command} /charts/{os.path.basename(chart_location)}"
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

        output = client.containers.run(verifier_image,docker_command,stdin_open=True,tty=True,stdout=True,remove=True,volumes=docker_volumes,environment=docker_environment)

    except docker.errors.ContainerError as exc:
        return f"FAIL: docker.errors.ContainerError: {exc.args}"
    except docker.errors.ImageNotFound as exc:
        return f"FAIL: docker.errors.ImageNotFound: {exc.args}"
    except docker.errors.APIError as exc:
        return f"FAIL: docker.errors.APIError: {exc.args}"

    if not output:
        return f"FAIL: no report produced : {docker_command}"

    return output.decode("utf-8")

def run_report_docker_image(verifier_image_name,verifier_image_tag,profile_type, report_location):

    verifier_image = f"{verifier_image_name}:{verifier_image_tag}"

    os.environ["VERIFIER_IMAGE"] = verifier_image
    docker_command = f"report all /reports/"+os.path.basename(report_location)

    if profile_type:
        docker_command = f"{docker_command} --set profile.vendorType={profile_type}"

    print(f'docker command: {docker_command}')

    try:
        client = docker.from_env()
        report_directory = os.path.dirname(os.path.abspath(report_location))
        output = client.containers.run(verifier_image,docker_command,stdin_open=True,tty=True,stdout=True,remove=True,volumes={report_directory: {'bind': '/reports/', 'mode': 'rw'}})
    except docker.errors.ContainerError as exc:
        return f"FAIL: docker.errors.ContainerError: {exc.args}"
    except docker.errors.ImageNotFound as exc:
        return f"FAIL: docker.errors.ImageNotFound: {exc.args}"
    except docker.errors.APIError as exc:
        return f"FAIL: docker.errors.APIError: {exc.args}"

    if not output:
        return f"FAIL: no report produced : {docker_command}"

    return output.decode("utf-8")

def run_version_tarball_image(tarball_name):
    tar = tarfile.open(tarball_name, "r:gz")
    tar.extractall(path="./test_verifier")
    out = subprocess.run(["./test_verifier/chart-verifier","version", "--as-data"],capture_output=True)
    return normalize_version(out.stdout.decode("utf-8"))

def normalize_version(version):
    """Extract normalized version from JSON data

    Parameters:
    version (string): the version and commit ID in JSON

    Returns:
    string: a normalized semver like 0.0.0.
    """
    print(f'version input to normalize_version function is: {version}')
    version_dict = json.loads(version)
    return version_dict["version"]

def run_version_docker_image(verifier_image_name,verifier_image_tag):
    """Run chart verifier's version command using the Docker image."""
    verifier_image = f"{verifier_image_name}:{verifier_image_tag}"
    os.environ["VERIFIER_IMAGE"] = verifier_image
    try:
        client = docker.from_env()
        output = client.containers.run(verifier_image,"version --as-data",stdin_open=True,tty=True,stdout=True,remove=True)
    except docker.errors.ContainerError as exc:
        return f"FAIL: docker.errors.ContainerError: {exc.args}"
    except docker.errors.ImageNotFound as exc:
        return f"FAIL: docker.errors.ImageNotFound: {exc.args}"
    except docker.errors.APIError as exc:
        return f"FAIL: docker.errors.APIError: {exc.args}"

    if not output:
        return f"FAIL: did not receive output from the chart verifier version subcommand."

    return normalize_version(output.decode("utf-8"))

def run_version_podman_image(verifier_image_name,verifier_image_tag):
    """Run chart verifier's version command in Podman."""
    out = subprocess.run(["podman", "run", "--rm", f"{verifier_image_name}:{verifier_image_tag}", "version", "--as-data"], capture_output=True)
    return normalize_version(out.stdout.decode("utf-8"))

def run_verify_tarball_image(tarball_name,profile_type, chart_location,pgp_key_location=None):
    print(f"Run tarball image from {tarball_name}")

    tar = tarfile.open(tarball_name, "r:gz")

    tar.extractall(path="./test_verifier")

    if pgp_key_location:
        out = subprocess.run(["./test_verifier/chart-verifier","verify","--set",f"profile.vendorType={profile_type}","--pgp-public-key",pgp_key_location,chart_location],capture_output=True)
    else:
        out = subprocess.run(["./test_verifier/chart-verifier","verify","--set",f"profile.vendorType={profile_type}",chart_location],capture_output=True)

    return out.stdout.decode("utf-8")

def run_report_tarball_image(tarball_name,profile_type, chart_location):
    print(f"Run tarball image from {tarball_name}")

    tar = tarfile.open(tarball_name, "r:gz")

    tar.extractall(path="./test_verifier")

    out = subprocess.run(["./test_verifier/chart-verifier","report","all","--set",f"profile.vendorType={profile_type}",chart_location],capture_output=True)

    return out.stdout.decode("utf-8")

def run_verify_podman_image(verifier_image_name,verifier_image_tag,profile_type, chart_location,pgp_key_location=None):

    print(f"Run podman image - {verifier_image_name}:{verifier_image_tag}")
    kubeconfig = os.environ.get("KUBECONFIG")
    if not kubeconfig:
        return "FAIL: missing KUBECONFIG environment variable"

    if chart_location.startswith('http:/') or chart_location.startswith('https:/'):
        if pgp_location:
            out = subprocess.run(["podman", "run", "-v", f"{kubeconfig}:/kubeconfig:z", "-e", "KUBECONFIG=/kubeconfig", "--rm",
                                  f"{verifier_image_name}:{verifier_image_tag}", "verify", "--set", f"profile.vendortype={profile_type}","--pgp-public-key",public_key_location,chart_location], capture_output=True)
        else:
            out = subprocess.run(["podman", "run", "-v", f"{kubeconfig}:/kubeconfig:z", "-e", "KUBECONFIG=/kubeconfig", "--rm",
                          f"{verifier_image_name}:{verifier_image_tag}", "verify", "--set", f"profile.vendortype={profile_type}", chart_location], capture_output=True)
    else:
        chart_directory = os.path.dirname(os.path.abspath(chart_location))
        chart_name = os.path.basename(os.path.abspath(chart_location))
        if pgp_key_location:
            pgp_key_name = os.path.basename(os.path.abspath(pgp_key_location))
            out = subprocess.run(["podman", "run", "-v", f"{chart_directory}:/charts:z", "-v", f"{kubeconfig}:/kubeconfig:z", "-e", "KUBECONFIG=/kubeconfig", "--rm",
                                  f"{verifier_image_name}:{verifier_image_tag}", "verify", "--set", f"profile.vendortype={profile_type}","--pgp-public-key",f"/charts/{pgp_key_name}",f"/charts/{chart_name}"], capture_output=True)
        else:
            out = subprocess.run(["podman", "run", "-v", f"{chart_directory}:/charts:z", "-v", f"{kubeconfig}:/kubeconfig:z", "-e", "KUBECONFIG=/kubeconfig", "--rm",
                              f"{verifier_image_name}:{verifier_image_tag}", "verify", "--set", f"profile.vendortype={profile_type}", f"/charts/{chart_name}"], capture_output=True)

    return out.stdout.decode("utf-8")

def run_report_podman_image(verifier_image_name,verifier_image_tag,profile_type, report_location):

    print(f"Run podman image - {verifier_image_name}:{verifier_image_tag}")

    report_directory = os.path.dirname(os.path.abspath(report_location))
    report_name = os.path.basename(os.path.abspath(report_location))
    out = subprocess.run(["podman", "run", "-v", f"{report_directory}:/reports:z", "--rm",
                        f"{verifier_image_name}:{verifier_image_tag}", "report", "all", "--set", f"profile.vendortype={profile_type}", f"/reports/{report_name}"], capture_output=True)

    return out.stdout.decode("utf-8")

@then("I should see the report-info from the report for the signed chart matching the expected report-info")
def signed_chart_report(run_signed_verify, profile_type, report_info_location, image_type, verifier_version):
    check_report(run_signed_verify, profile_type, report_info_location, image_type, verifier_version)


@then("I should see the report-info from the generated report matching the expected report-info")
def chart_report(run_verify, profile_type, report_info_location, image_type, verifier_version):
    check_report(run_verify, profile_type, report_info_location, image_type, verifier_version)

def check_report(verify_result, profile_type, report_info_location, image_type, verifier_version):

    if verify_result.startswith("FAIL"):
        pytest.fail(f'FAIL some tests failed: {verify_result}')

    print(f"Report data:\n{verify_result}\ndone")

    report_data = yaml.load(verify_result, Loader=Loader)


    test_passed = True

    report_verifier_version = report_data['metadata']['tool']['verifier-version']
    if report_verifier_version != verifier_version:
        print(f"FAIL: verifier-version found in report does not match tool version. Expected {verifier_version}, but report has {report_verifier_version}")
        test_passed = False

    report_vendor_type = report_data["metadata"]["tool"]["profile"]["VendorType"]
    if report_vendor_type != profile_type:
        print(f"FAIL: profiles do not match. Expected {profile_type}, but report has {report_vendor_type}")
        test_passed = False

    chart_name =  report_data["metadata"]["chart"]["name"]
    chart_version = report_data["metadata"]["chart"]["version"]

    report_name = f'{profile_type}-{chart_name}-{chart_version}-report.yaml'
    report_dir = os.path.join(os.getcwd(),"test-reports")
    report_path = os.path.join(report_dir,report_name)

    if not os.path.isdir(report_dir):
        os.makedirs(report_dir)

    print(f'Report path : {report_path}')

    with open(report_path, "w") as fd:
        fd.write(verify_result)

    expected_reports_file = open(report_info_location,)
    expected_reports = json.load(expected_reports_file)

    test_report_data = run_report(image_type, profile_type, report_path)
    print(f"test_report_data : \n{test_report_data}")
    test_report = json.loads(test_report_data)

    results_diff = DeepDiff(expected_reports[report_info.REPORT_RESULTS],test_report[report_info.REPORT_RESULTS],ignore_order=True)
    if results_diff:
        print(f"difference found in results : {results_diff}")
        test_passed = False

    expected_annotations = {}
    for report_annotation in expected_reports[report_info.REPORT_ANNOTATIONS]:
        expected_annotations[report_annotation["name"]] = report_annotation["value"]

    tested_annotations = {}
    for report_annotation in test_report[report_info.REPORT_ANNOTATIONS]:
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

    digests_diff = DeepDiff(expected_reports[report_info.REPORT_DIGESTS],test_report[report_info.REPORT_DIGESTS],ignore_order=True)
    if digests_diff:
        print(f"difference found in digests : {digests_diff}")
        test_passed = False

    differing_metadata = {"chart-uri"}
    expected_metadata = expected_reports[report_info.REPORT_METADATA]
    test_metadata = test_report[report_info.REPORT_METADATA]
    for metadata in expected_metadata.keys():
        if not metadata in differing_metadata:
            metadata_diff = DeepDiff(expected_metadata[metadata],test_metadata[metadata],ignore_order=True)
            if metadata_diff:
                test_passed = False
                print(f"difference found in {metadata} metadata : {metadata_diff}")

    if not test_passed:
        pytest.fail('FAIL differences found in reports')

