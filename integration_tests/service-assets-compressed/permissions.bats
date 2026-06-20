#!/usr/bin/env bats

@test "forte installs service assets with correct permissions" {
    [ "$(stat -c '%a' /srv/assets)"                  = "750" ]
    [ "$(stat -c '%a' /srv/assets/fake_app/app.css)" = "640" ]
}
