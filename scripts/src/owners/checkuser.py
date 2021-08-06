import re
import argparse
import requests
import os
import sys
import yaml
try:
    from yaml import CLoader as Loader, CDumper as Dumper
except ImportError:
    from yaml import Loader, Dumper

def verify_user(username):
    print(f"[INFO] Verify user. {username}")
    owners_path = "OWNERS"
    if not os.path.exists(owners_path):
        print(f"[ERROR] {owners_path} file does not exist.")
    else:
        data = open(owners_path).read()
        out = yaml.load(data, Loader=Loader)
        if username in out["approvers"]:
            print(f"[INFO] {username} authorized")
            return True
        else:
            print(f"[ERROR] {username} not auhtorized")
    return False

def check_for_restricted_file(api_url):
    files_api_url = f'{api_url}/files'
    headers = {'Accept': 'application/vnd.github.v3+json'}
    pattern_owners = re.compile(r"OWNERS")
    pattern_versionfile = re.compile(r"cmd/release/release_info.json")
    pattern_thisfile = re.compile(r"scripts/src/owners/checkuser.py")
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
           if pattern_versionfile.match(filename) or pattern_owners.match(filename) or pattern_thisfile.match(filename):
               print(f"[INFO] restricted file found: {filename}")
               return True
 
    return False


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-a", "--api-url", dest="api_url", type=str, required=False,
                        help="API URL for the pull request")
    parser.add_argument("-u", "--user", dest="username", type=str, required=False,
                        help="check if the user can run tests")
    args = parser.parse_args()

    if check_for_restricted_file(args.api_url):
        if verify_user(args.username):
            print(f"[INFO] {args.username} is authorized to modify all files in the PR")
        else:
            print(f"[INFO] {args.username} is not authorized to modify all files in the PR")
            sys.exit(1)
    else:
        print(f"[INFO] no restricted files found in the PR")