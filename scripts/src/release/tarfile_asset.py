import argparse
import os
import tarfile


tar_content_files = [ {"name": "out/chart-verifier", "arc_name": "chart-verifier"} ]


def create(release):

    tgz_name = f"chart-verifier-{release}.tgz"
    print(f'::set-output name=tarball_base_name::{tgz_name}')

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
    print(f'::set-output name=tarball_full_name::{tarfile}')

