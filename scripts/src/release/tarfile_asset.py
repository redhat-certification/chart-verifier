import argparse
import os
import tarfile
from utils import utils

tar_content_files = [ {"name": "out/chart-verifier", "arc_name": "chart-verifier"} ]


def create(release):

    tgz_name = f"chart-verifier-{release}.tgz"
    utils.add_output("tarball_base_name",tgz_name)

    if os.path.exists(tgz_name):
        os.remove(tgz_name)

    with tarfile.open(tgz_name, "x:gz") as tar:
        for tar_content_file in tar_content_files:
            tar.add(os.path.join(os.getcwd(),tar_content_file["name"]),arcname=tar_content_file["arc_name"])


    return os.path.join(os.getcwd(),tgz_name)

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-r", "--release", dest="release", type=str, required=True,
                        help="Release name for the tar file")

    args = parser.parse_args()
    tarfile = create(args.release)
    print(f'[INFO] Verifier tarball created : {tarfile}.')
    utils.add_output("tarball_full_name",tarfile)

