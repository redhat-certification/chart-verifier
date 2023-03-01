"""
Used by a github action to build a test a chart verifier docker image based on a pull request

parameters:
    ---image-name : Name to be used for the chart verifier image.
    --sha-value : Sha value to be used for the image.
    --verifier-version : Version of the chart-verifier to test.

results:
    exit code 1 if the image fails to build or a test fails.
"""

import argparse
import docker
import os
import sys
import yaml
import json
try:
    from yaml import CLoader as Loader, CDumper as Dumper
except ImportError:
    from yaml import Loader, Dumper

sys.path.append('./scripts/src/')
from report import report_info
from utils import utils

def build_image(image_id):
    print(f"Build Image : {image_id}")

    cwd = os.getcwd()

    # Print the current working directory
    print(f"Current working directory: {cwd}")

    client = docker.from_env()

    try:
        image = client.images.build(path="./",tag=image_id)
        print("images:",image)
    except docker.errors.BuildError as buildError:
        print("docker build error")
        for line in buildError.build_log:
            if 'stream' in line:
                logger.error(line['stream'].strip())
        return False
    except  docker.errors.APIError:
        print ("docker API error")
        return False

    return True

def test_image(image_id,chart,verifier_version):

    docker_command = "verify -l " + chart["url"]

    set_values = ""
    vendor_type = ""
    profile_version = ""
    if "vendorType" in chart["metadata"]:
        vendor_type = chart["metadata"]["vendorType"]
        set_values = "profile.vendortype=%s" % vendor_type
    if "profileVersion" in chart["metadata"]:
        profile_version = chart["metadata"]["profileVersion"]
        if set_values:
            set_values = "%s,profile.version=%s" % (set_values,profile_version)
        else:
            set_values = "profile.version=%s" % profile_version

    if set_values:
        docker_command = "%s --set %s" % (docker_command, set_values)

    client = docker.from_env()
    out = client.containers.run(image_id,docker_command,stdin_open=True,tty=True,stderr=True)
    report = yaml.load(out, Loader=Loader)
    report_path = "banddreport.yaml"

    if verifier_version and verifier_version != report["metadata"]["tool"]["verifier-version"]:
        print(f'[ERROR] Chart verifier report version {report["metadata"]["tool"]["verifier-version"]} does not match  expected version: {verifier_version}')
        return False

    docker_command = "version"
    # sample output: v1.0.0
    out = client.containers.run(image_id, docker_command, stdin_open=True, tty=True, stderr=True)
    if not out or out[1:] != verifier_version:
        print(f"[ERROR] 'chart-verifier version' output {out} does not match expected version: {verifier_version}")

    print("[INFO] report:\n", report)
    with open(report_path, "w") as fd:
        yaml.dump(report,fd)

    results = report_info.get_report_results(report_path,vendor_type,profile_version)

    expectedPassed = int(chart["results"]["passed"])
    expectedFailed = int(chart["results"]["failed"])

    if expectedFailed != results["failed"] or expectedPassed != results["passed"]:
        print("[ERROR] Chart verifier report includes unexpected results:")
        print(f'- Number of checks passed expected : {expectedPassed}, got {results["passed"]}')
        print(f'- Number of checks failed expected : {expectedFailed}, got {results["failed"]}')
        return False
    else:
        print(f'[PASS] Chart result validated : {chart["url"]}')

    return True


def main():

    utils.add_output("result","failure")

    parser = argparse.ArgumentParser()
    parser.add_argument("-i", "--image-name", dest="image_name", type=str, required=True,
                        help="Name of the chart verifier image")
    parser.add_argument("-s", "--sha-value", dest="sha_value", type=str, required=True,
                        help="Image sha value to test")
    parser.add_argument("-v", "--verifier-version", dest="verifier_version", type=str, required=False,
                        help="New version of chart verifier")
    parser.add_argument("-b","--build-only",dest="build_only",type=str, required=False,
                        help="set to any value to build only and not test")


    args = parser.parse_args()

    image_id = f"{args.image_name}:{args.sha_value}"

    if build_image(image_id):

        utils.add_output("verifier-image-tag",args.sha_value)

        if not args.build_only:

            chart = {"url" : "https://github.com/redhat-certification/chart-verifier/blob/main/pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz?raw=true",
                "results":{"passed":"10","failed":"1"},
                "metadata":{"vendorType":"partner","profileVersion":"v1.0"}}

            os.environ["VERIFIER_IMAGE"] = image_id

            if not test_image(image_id,chart,args.verifier_version):
                sys.exit(1)

    else:
        sys.exit(1)
