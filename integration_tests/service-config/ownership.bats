#!/usr/bin/env bats

@test "forte installs service config with correct ownership" {
    [ "$(stat -c '%U:%G' /etc/fake_app)"                     = "root:www-data" ]
    [ "$(stat -c '%U:%G' /etc/fake_app/app_info.php)"        = "root:www-data" ]
}
