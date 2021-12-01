
import os
import sys
import docker
import json

from deepdiff import DeepDiff

REPORT_ANNOTATIONS = "annotations"
REPORT_RESULTS = "results"
REPORT_DIGESTS = "digests"
REPORT_METADATA = "metadata"
ALL_REPORTS = "all"

def _get_report_info(report_path, info_type, profile_type, profile_version):

    docker_command = "report " + info_type + " /charts/"+os.path.basename(report_path)

    set_values = ""
    if profile_type:
        set_values = "profile.vendortype=%s" % profile_type
    if profile_version:
        if set_values:
            set_values = "%s,profile.version=%s" % (set_values,profile_version)
        else:
            set_values = "profile.version=%s" % profile_version

    if set_values:
        docker_command = "%s --set %s" % (docker_command, set_values)

    client = docker.from_env()
    report_directory = os.path.dirname(os.path.abspath(report_path))
    output = client.containers.run(os.environ.get("VERIFIER_IMAGE"),docker_command,stdin_open=True,tty=True,stderr=True,volumes={report_directory: {'bind': '/charts/', 'mode': 'rw'}})
    report_out = json.loads(output)

    if info_type == ALL_REPORTS:
        return report_out

    if not info_type in report_out:
        print(f"Error extracting {info_type} from the report:", report_out.strip())
        sys.exit(1)

    if info_type == REPORT_ANNOTATIONS:
        annotations = {}
        for report_annotation in report_out[REPORT_ANNOTATIONS]:
            annotations[report_annotation["name"]] = report_annotation["value"]

        return annotations

    return report_out[info_type]


def get_report_annotations(report_path):
    annotations = _get_report_info(report_path,REPORT_ANNOTATIONS,"","")
    print("[INFO] report annotations : %s" % annotations)
    return annotations

def get_report_results(report_path, profile_type, profile_version):
    results = _get_report_info(report_path,REPORT_RESULTS,profile_type,profile_version)
    print("[INFO] report results : %s" % results)
    results["failed"] = int(results["failed"])
    results["passed"] = int(results["passed"])
    return results
    
def get_report_digests(report_path):
    digests = _get_report_info(report_path,REPORT_DIGESTS,"","")
    print("[INFO] report digests : %s" % digests)
    return digests

def get_report_metadata(report_path):
    metadata = _get_report_info(report_path,REPORT_METADATA,"","")
    print("[INFO] report metadata : %s" % metadata)
    return metadata

def get_report_chart_url(report_path):
     metadata = _get_report_info(report_path,REPORT_METADATA,"","")
     print("[INFO] report chart-uri : %s" % metadata["chart-uri"])
     return metadata["chart-uri"]

def get_report_chart(report_path):
     metadata = _get_report_info(report_path,REPORT_METADATA,"","")
     print("[INFO] report chart : %s" % metadata["chart"])
     return metadata["chart"]

def get_all_reports(report_path,profile_type,profile_version):
    all_reports = _get_report_info(report_path,ALL_REPORTS,profile_type,profile_version)
    print(f"[INFO] all reports : {all_reports}")
    return all_reports






