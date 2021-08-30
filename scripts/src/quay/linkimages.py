"""
Used by a github action to link a newly published chart verifier image to a specified tag in quay
Can also be called from the command line.

parameters:
    ---link-tag : tag to link the image to, default is "test"
    --verifier-version : the version of the new chart verifier release

notes:
    It is anticipated that the github action will invoke this function just before the new image is available
    in quay. As a result the function will look for the image every 15 seconds for up to 15 minutes.

    To invoke from the command line, invoke from the root directory of the repository,for example:
          python3 scripts/src/quay/linkimages.py --verifier-version=1.2.0.

    If the verifier version image is already linked to the specified tag no action is taken.

    The default of link-tag is set to "test" to avoid accidental updates to latest.

    The auth token to enable the script to link tags in quay must be set in an environment variable QUAY_AUTH_TOKEN.
    For a github action this must be set as a repository secret.
    For information on auth token required see:
    https://access.redhat.com/documentation/en-us/red_hat_quay/3/html/red_hat_quay_api_guide/using_the_red_hat_quay_api

results:
    A message indicating the outcome.
    exit code 1 if version image was not found.

"""

import requests
import json
import sys
import argparse
import os
from retry import retry

sys.path.append('./scripts/src/')
from release import releasechecker

# Quay API docs: https://access.redhat.com/documentation/en-us/red_hat_quay/3/html/red_hat_quay_api_guide/index
# Quay Swagger API: https://docs.quay.io/api/swagger/
#

DEFAULT_TAG_TO_LINK= "test"
QUAY_TAG_PAGE_URL = 'https://quay.io/api/v1/repository/redhat-certification/chart-verifier/tag/'

# try every 15 seconds for 15 minutes
@retry(Exception,tries=60, delay=15)
def get_image_id(tag_value, do_retry):

    print(f"[INFO] look for tag : {tag_value}, retry : {do_retry}")
    tag_url = 'https://quay.io/api/v1/repository/redhat-certification/chart-verifier/tag/'

    get_params = {'onlyActiveTags' : 'true','specificTag' : tag_value}

    response = requests.get(tag_url,params=get_params)

    image_id = ""
    if response.status_code not in [200,201]:
        print(f"[Error] Error getting tags from quay : status_code={response.status_code}")
    else:
        tags = json.loads(response.text)
        print("[INFO] loaded the tags")
        for tag in tags["tags"]:
            if tag['name'] == tag_value:
                image_id = tag['image_id']
                print(f"[INFO] Found tag {tag_value}. image_id : {image_id}")
                break
            else:
                print(f"[INFO] ignore tag {tag['name']}")

        if not image_id and do_retry:
            print(f"[INFO] {tag_value} not found. Retry!")
            raise Exception(f"Image {tag_value} not found")

    return image_id

def link_image(image_to_link, tag_value):

    print(f"[INFO] Update {tag_value} to point to {image_to_link}")
    auth_token = os.environ.get('QUAY_AUTH_TOKEN')
    if not auth_token:
        print("[ERROR] repository secret QUAY_AUTH_TOKEN not set")
        return False

    quay_token = f"Bearer {auth_token}"
    put_header = {'content-type': 'application/json','Authorization': quay_token}

    put_url = f"{QUAY_TAG_PAGE_URL}{tag_value}"
    put_data = {'image': image_to_link}
    put_out = requests.put(put_url,data=json.dumps(put_data), headers=put_header)
    print(f"[INFO] Update link response code : {put_out.status_code}")
    print(f"[INFO] Update link response : {put_out.text}")

    return put_out.status_code in [200,201]


def main():

    parser = argparse.ArgumentParser()
    parser.add_argument("-t", "--link-tag", dest="link_tag", type=str, required=False, default=DEFAULT_TAG_TO_LINK,
                        help="Tag image should be linked to (default: test")
    parser.add_argument("-v", "--verifier-version", dest="verifier_version", type=str, required=False,
                        help="New version of chart verifier")
    args = parser.parse_args()

    new_tag = args.verifier_version
    if not args.verifier_version:
        version_info = releasechecker.get_version_info()
        new_tag = version_info["version"]

    try:
        new_image_id = get_image_id(new_tag, True)
        if not new_image_id:
            print(f"[ERROR] Failed find new Image : {new_tag}")
            sys.exit(1)
        tag_image_id = get_image_id(args.link_tag, False)
        if tag_image_id != new_image_id:
            if link_image(new_image_id, args.link_tag):
                print(f"[INFO] PASS {args.link_tag} linked to {new_tag}")
                return
            else:
                print(f"[ERROR] Failed to link tags")
                sys.exit(1)
        else:
            print(f"[INFO] Tag {args.link_tag} is current")
            return
    except Exception as inst:
        print(f"[WARNING] {inst.args}")
        sys.exit(1)

    return

if __name__ == "__main__":
    main()
