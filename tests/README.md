# Chart Verifier python behave bdd tests

## Regular chart tests

Tests are basic and test good charts only:
1. chart src ```tests/charts/psql-service/0.1.8/src```
2. chart tgz ```tests/charts/psql-service/0.1.9/psql-service-0.1.9.tgz```

Each chart is verified as partner, redhat and community. A report summary, of the verifier report created, is obtained for each chart and the content is compared to expected content. 

This amounts to 6 tests. These 6 tests are then run on 3 image types:
- chart-verifier docker image
- chart-verifier podman image
- chart-verifier binary image

As a result there are 18 tests.

## Signed chart tests

Signed chart tests also include only a test of a valid signed chart:
   - ```tests/charts/psql-service/0.1.11/psql-service-0.1.9.tgz```

To be seen as signed a helm providence file is included:
   - ```tests/charts/psql-service/0.1.11/psql-service-0.1.9.tgz.prov```

And to enable verification of the chart a public key file is included which contains the public key of the secret key used to sign that chart:
   - ```tests/charts/psql-service/0.1.11/psql-service-0.1.9.tgz.prov```

The chart is verified as partner and redhat. A report summary, of the verifier report created, is obtained for the chart and the content is compared to expected content.

This amounts to 2 tests. These 2 tests are then run on 3 image types:
- chart-verifier docker image
- chart-verifier podman image
- chart-verifier binary image

As a result there are 6 tests.

### Signing the chart

The chart is signed using helm cli and a secret key. In this initial version the secret key used was one from Martin Mulholland. As a result the public key checked in for the test is also from Martin.

In the event the chart has to be updated, or a new chart added, the creator or updater of the chart can sign it use their own secret key, and create a copy of their public key to include with the test. 

This is not ideal, we will investigate generating secret and public keys as art of the test using a bot id. 

This is a useful link for [gpg cheat sheet](http://irtfweb.ifa.hawaii.edu/~lockhart/gpg/).