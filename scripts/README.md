##Scripts for chart verifier workflows


- ```release/versionchecker.py```
    - used to determine if a PR contains a version update. 
-  ```buildandtest/buildandtest.py```
    - used to build a docker image and then test created image.
- ```checkautomerge/checkautomerge.py```
    - loops waiting for a PR to merge
    - exact copy of same script from chart repo
- ```report/rrport-info.py```
    - used to generate of report of a chart verifier verify report.
    - exact copy of same script from chart repo    
    