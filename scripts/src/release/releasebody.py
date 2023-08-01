import sys

sys.path.append("./scripts/src/")
from release import releasechecker  # noqa E402


def get_release_body(version, image_name, release_info):
    """Generate the body of the GitHub release"""
    body = f"Chart verifier version {version} <br><br>"
    body += f"Docker Image:<br>- {image_name}:{version}<br><br>"
    body += "This version includes:<br>"
    for info in release_info:
        if info.startswith("<"):
            body += info
        else:
            body += f"- {info}<br>"
    return body


def main():
    version_info = releasechecker.get_version_info()
    release_body = get_release_body(
        version_info["version"],
        version_info["quay-image"],
        version_info["release-info"],
    )
    print(release_body)
