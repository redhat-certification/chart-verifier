import os

def add_output(name,value):
    with open(os.environ['GITHUB_OUTPUT'],'a') as fh:
        print(f'{name}={value}',file=fh)
