import os

def add_output(key,value):
    """This function prints the key/value pair to stdout in a "key=value" format.

    If called from a GitHub workflow, it also sets an output parameter.
    """

    print(f'{key}={value}')

    if "GITHUB_OUTPUT" in os.environ:
        with open(os.environ['GITHUB_OUTPUT'],'a') as fh:
            print(f'{key}={value}',file=fh)
