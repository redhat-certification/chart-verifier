"""
Used by a github action
1. To determine if the contents of pull request contain only the file which contains the chart verifier release.
2. To determine if the release has been updated.

parameters:
    --api-url : API URL for the pull request.
    --version : user to be checked for authority to modify release files in a PR.

results:
    if --api-url is specified, output variables are set:
        PR_version : The chart verifier version read from the version file from the PR.
        PR_release_image : The name of the image to be used for the release.
        PR_release_info : Information about the release content.
        PR_includes_release : Set to true if the PR contains the version file.
        PR_release_body : Body of text to be used to describe the release.
    if --version only is specified, output variables are set:
        updated : set to true if the version specified is later than the version in the version file
                  from the main branch.
    if neither parameters are specified, output variables are set:
        PR_version : The chart verifier version read from the version file from main branch.
        PR_release_image : The name of the image from the version file from main branch.
"""

import re
import argparse
import json
import requests
import semver

VERSION_FILE = "cmd/release/release_info.json"

def check_if_only_version_file_is_modified(api_url):
    # api_url https://api.github.com/repos/<organization-name>/<repository-name>/pulls/<pr_number>

    files_api_url = f'{api_url}/files'
    headers = {'Accept': 'application/vnd.github.v3+json'}
    pattern_versionfile = re.compile(r"cmd/release/release_info.json")
    page_number = 1
    max_page_size,page_size = 100,100


    version_file_found = False
    while (page_size == max_page_size):

        files_api_query = f'{files_api_url}?per_page={page_size}&page={page_number}'
        r = requests.get(files_api_query,headers=headers)
        files = r.json()
        page_size = len(files)
        page_number += 1

        for f in files:
            filename = f["filename"]
            if pattern_versionfile.match(filename):
                version_file_found = True
            else:
                return False

    return version_file_found

def make_release_body(version, image_name, release_info):
    body = f"Chart verifier version {version} <br><br>Docker Image:<br>- {image_name}:{version}<br><br>"
    body += "This version includes:<br>"
    for info in release_info:
        if info.startswith("<"):
            body += info
        else:
            body += f"- {info}<br>"

    print(f"[INFO] Release body: {body}")
    print(f"::set-output name=PR_release_body::{body}")

def get_version_info():
    data = {}
    with open(VERSION_FILE) as json_file:
        data = json.load(json_file)
    return data

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-a", "--api-url", dest="api_url", type=str, required=False,
                        help="API URL for the pull request")
    parser.add_argument("-v", "--version", dest="version", type=str, required=False,
                        help="Version to compare")

    args = parser.parse_args()
    if args.api_url and check_if_only_version_file_is_modified(args.api_url):
        ## should be on PR branch
        version_info = get_version_info()
        print(f'[INFO] Release found in PR files : {version_info["version"]}.')
        print(f'::set-output name=PR_version::{version_info["version"]}')
        print(f'::set-output name=PR_release_image::{version_info["quay-image"]}')
        print(f'::set-output name=PR_release_info::{version_info["release-info"]}')
        print(f'::set-output name=PR_includes_release::true')
        make_release_body(version_info["version"],version_info["quay-image"],version_info["release-info"])
    else:
        version_info = get_version_info()
        if args.version:
            # should be on main branch
            if semver.compare(args.version,version_info["version"]) > 0 :
                print(f'[INFO] Release {args.version} found in PR files is newer than: {version_info["version"]}.')
                print("::set-output name=updated::true")
            else:
                print(f'[INFO] Release found in PR files is not new  : {version_info["version"]}.')
        else:
            print(f'::set-output name=PR_version::{version_info["version"]}')
            print(f'::set-output name=PR_release_image::{version_info["quay-image"]}')
            print("[INFO] PR contains non-release files.")
