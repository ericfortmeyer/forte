#!/usr/bin/env bats

setup_file() {
    cp -r ./integration_tests/testdata/fake_app-assets /tmp/fake_app-assets

    forte deploy fake_app "www-data"
}

@test "forte performs FHS compliant service assets installations" {
    [ -d /srv/assets/fake_app ] # See FHS-3.0 § 3.17
}
