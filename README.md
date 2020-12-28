

- [ ] Command line scaffold.
    - `helmcertifier --only check1,check2 ./chart.tgz`
    - [x] `--only check1,check2` to perform only the checks specified in the comma-separated list.
    - [x] `--except check1,check2` to perform all checks *except* those specified in the comma-separated list.
    - [x] `--uri` representing a value Helm understands as a *chart URI*; required. 
- [ ] Business logic package.
    - Configuration as input, check result as output.
 