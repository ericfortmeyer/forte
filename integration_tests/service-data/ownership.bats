#!/usr/bin/env bats

@test "forte installs service data with correct ownership" {
    [ "$(stat -c '%U:%G' /srv/fake_app)"                     = "root:www-data" ]
    [ "$(stat -c '%U:%G' /srv/fake_app/public)"              = "root:www-data" ]
    [ "$(stat -c '%U:%G' /srv/fake_app/src)"                 = "root:www-data" ]
    [ "$(stat -c '%U:%G' /srv/fake_app/public/index.php)"    = "root:www-data" ]
    [ "$(stat -c '%U:%G' /srv/fake_app/src/ItemService.php)" = "root:www-data" ]
}
