#!/usr/bin/env bats

@test "forte installs service config with correct permissions" {
    [ "$(stat -c '%a' /etc/fake_app)"                     = "750" ]
    [ "$(stat -c '%a' /etc/fake_app/app_info.php)"        = "640" ]
}
