#!/usr/bin/env bats

@test "forte installs service data with correct permissions" {
    [ "$(stat -c '%a' /srv/fake_app)"                     = "750" ]
    [ "$(stat -c '%a' /srv/fake_app/public)"              = "750" ]
    [ "$(stat -c '%a' /srv/fake_app/src)"                 = "750" ]
    [ "$(stat -c '%a' /srv/fake_app/public/index.php)"    = "640" ]
    [ "$(stat -c '%a' /srv/fake_app/src/ItemService.php)" = "640" ]
}
