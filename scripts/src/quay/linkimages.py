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

defaultLinkTag= "test"
tagUrl = 'https://quay.io/api/v1/repository/redhat-certification/chart-verifier/tag/'

# try every 15 seconds for 15 minutes
@retry(Exception,tries=10, delay=5)
def getImageId(tagValue,doRetry):

    print(f"[INFO] look for tag : {tagValue}, retry : {doRetry}")
    tagUrl = 'https://quay.io/api/v1/repository/redhat-certification/chart-verifier/tag/'

    getParams = {'onlyActiveTags' : 'true','specificTag' : tagValue }

    response = requests.get(tagUrl,params=getParams)

    imageId = ""
    if response.status_code > 201:
        print(f"[Error] Error getting tags from quay : status_code={response.status_code}")
    else:
        tags = json.loads(response.text)
        print("[INFO] loaded the tags")
        for tag in tags["tags"]:
            if tag['name'] == tagValue:
                imageId = tag['image_id']
                print(f"[INFO] Found tag {tagValue}. image_id : {imageId}")
                break
            else:
                print(f"[INFO] ignore tag {tag['name']}")

        if not imageId and doRetry:
            print(f"[INFO] {tagValue} not found. Retry!")
            raise Exception(f"Image {tagValue} not found")

    return imageId

def linkImage(linkImage,linkTag):

    print(f"[INFO] Update {linkTag} to point to {linkImage}")
    auth_token = f"Bearer {os.environ.get('QUAY_AUTH_TOKEN')}"
    putHeader = {'content-type': 'application/json','Authorization': auth_token }

    puturl = tagUrl + linkTag
    putData = {'image': linkImage}
    putOut = requests.put(puturl,data=json.dumps(putData), headers=putHeader)
    print(f"[INFO] Update link response code : {putOut.status_code}")
    print(f"[INFO] Update link response : {putOut.text}")

    return putOut.status_code == 200 or putOut.status_code == 201


def main():

    parser = argparse.ArgumentParser()
    parser.add_argument("-t", "--link-tag", dest="link_tag", type=str, required=False, default=defaultLinkTag,
                        help="Tag image should be linked to (default: test")
    parser.add_argument("-v", "--verifier-version", dest="verifier_version", type=str, required=False,
                        help="New version of chart verifier")
    args = parser.parse_args()

    newTag = args.verifier_version
    if args.verifier_version is None:
        version_info = releasechecker.get_version_info()
        newTag = version_info["version"]

    try:
        newImageId = getImageId(newTag,True)
        tagImageId = getImageId(args.link_tag,False)
        if tagImageId != newImageId:
            if linkImage(newImageId,args.link_tag):
                print(f"[INFO] PASS {args.link_tag} linked to {newTag}")
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
