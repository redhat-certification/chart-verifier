import tarfile
import os


tar_content_files = [ {"name": "config", "arc_name": "config"},
            {"name": "out/chart-verifier", "arc_name": "chart-verifier"} ]


def create(release):

    tgz_name = f"chart-verifier-{release}.tgz"

    if os.path.exists(tgz_name):
        os.remove(tgz_name)

    with tarfile.open(tgz_name, "x:gz") as tar:
        for tar_content_file in tar_content_files:
            tar.add(os.path.join(os.getcwd(),tar_content_file["name"]),arcname=tar_content_file["arc_name"])


    return os.path.join(os.getcwd(),tgz_name)

