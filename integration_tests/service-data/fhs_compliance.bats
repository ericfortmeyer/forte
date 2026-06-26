#!/usr/bin/env bats

setup_file() {
    cp -r ./integration_tests/testdata/fake_app /tmp/fake_app

    forte deploy fake_app "www-data" || exit 1
}

@test "forte performs FHS compliant service data installations" {
    [ -d /srv/fake_app ] # See FHS-3.0 § 3.17
}
