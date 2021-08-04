import re
import argparse
import json
import requests
import semver

version_file = "cmd/release/release_info.json"

def check_if_version_file_is_modified(api_url):
    # api_url https://api.github.com/repos/<organization-name>/<repository-name>/pulls/<pr_number>

    files_api_url = f'{api_url}/files'
    headers = {'Accept': 'application/vnd.github.v3+json'}
    pattern_versionfile = re.compile(version_file)
    page_number = 1
    max_page_size,page_size = 100,100


    while (page_size == max_page_size):

        files_api_query = f'{files_api_url}?per_page={page_size}&page={page_number}'
        r = requests.get(files_api_query,headers=headers)
        files = r.json()
        page_size = len(files)
        page_number += 1


        for f in files:
            filename = f["filename"]
            if pattern_versionfile.match(filename):
                return True

    return False

def make_relase_body(version, image_name, release_info):
    body = f"Chart verifier version {version} <br><br>Docker Image:<br>- {image_name}<br><br>"
    body += "This version includes:<br>"
    for info in release_info:
        body += f"- {info}<br>"

    print(f"[INFO] Release body: {body}")
    print(f"::set-output name=PR_release_body::{body}")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-a", "--api-url", dest="api_url", type=str, required=False,
                        help="API URL for the pull request")
    parser.add_argument("-v", "--version", dest="version", type=str, required=False,
                        help="Version to compare")
    args = parser.parse_args()
    if args.api_url and check_if_version_file_is_modified(args.api_url):
        ## should be on PR branch
        version_info = json.loads(version_file)
        print(f'[INFO] Release found in PR files : {version_info["version"]}.')
        print(f'::set-output name=PR_version::{version_info["version"]}')
        print(f'::set-output name=PR_release_image::{version_info["quay-image"]}')
        print(f'::set-output name=PR_release_info::{version_info["release-info"]}')
    elif args.version:
        # should be on main branch
        version_info = json.loads(version_file)
        if semver.compare(args.version,version_info["version"]) > 0 :
            print(f'[INFO] Release {args.Verision} found in PR files is newer than: {version_info["version"]}.')
            print("::set-output name=updated::true")
        else:
            print(f'[INFO] Release found in PR files is not new  : {version_info["version"]}.')
    else:
        print("[INFO] No new release found in PR files.")


