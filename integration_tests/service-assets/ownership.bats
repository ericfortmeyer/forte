#!/usr/bin/env bats

@test "forte installs service assets with correct ownership" {
    [ "$(stat -c '%U:%G' /srv/assets/fake_app)"         = "root:www-data" ]
    [ "$(stat -c '%U:%G' /srv/assets/fake_app/app.css)" = "root:www-data" ]
}
