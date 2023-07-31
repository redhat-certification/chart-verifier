# Chart Verifier python behave bdd tests

## Regular chart tests

Tests are basic and test good charts only:
1. chart src ```tests/charts/psql-service/0.1.8/src```
2. chart tgz ```tests/charts/psql-service/0.1.9/psql-service-0.1.9.tgz```

Each chart is verified as partner, redhat and community. A report summary, of the verifier report created, is obtained for each chart and the content is compared to expected content. 

This amounts to 6 tests. These 6 tests are then run on 2 image types:
- chart-verifier podman image
- chart-verifier binary image

As a result there are 12 tests.

## Signed chart tests

Signed chart tests also include only a test of a valid signed chart:
   - ```tests/charts/psql-service/0.1.11/psql-service-0.1.9.tgz```

To be seen as signed a helm providence file is included:
   - ```tests/charts/psql-service/0.1.11/psql-service-0.1.9.tgz.prov```

And to enable verification of the chart a public key file is included which contains the public key of the secret key used to sign that chart:
   - ```tests/charts/psql-service/0.1.11/psql-service-0.1.9.tgz.prov```

The chart is verified as partner and redhat. A report summary, of the verifier report created, is obtained for the chart and the content is compared to expected content.

This amounts to 2 tests. These 2 tests are then run on 2 image types:
- chart-verifier podman image
- chart-verifier binary image

As a result there are 4 tests.

### Signing the chart

The signed chart tests have been signed with a key generated specifically for
these tests. When these are changed, a new keypair must be generated to use for
signing. The secret key can be thrown away. The private key can be thrown away.
The public key is all that's required for these tests to complete, and this key
is not to be used for anything else.

TODO: Generate a workflow that does this automagically in a container, etc.

This is not ideal, we will investigate generating secret and public keys as art of the test using a bot id. 

This is a useful link for [gpg cheat sheet](http://irtfweb.ifa.hawaii.edu/~lockhart/gpg/).