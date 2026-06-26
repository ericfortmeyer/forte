#!/usr/bin/env bats

setup_file() {
    cp -r ./integration_tests/testdata/fake_app-config.tar.gz /tmp/fake_app-config.tar.gz

    forte deploy fake_app "www-data" || exit 1
}

@test "forte performs FHS compliant service config installations" {
    [ -d /etc/fake_app ] # See FHS-3.0 § 3.7
}
