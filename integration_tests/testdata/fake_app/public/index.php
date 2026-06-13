<?php

require_once __DIR__ . "/../src/ItemService.php";

$itemService = new ItemService();
$title = "A Fake App";
$bodyContent = <<<EOM
BLA BLA BLA BLA BLA
BLA BLA BLA BLA BLA
BLA BLA BLA BLA BLA
===================
BLA BLA BLA BLA BLA
BLA BLA BLA BLA BLA
BLA BLA BLA BLA BLA
===================
EOM;
?>
<!DOCTYPE html>
<head>
    <title><?= $title ?></title>
</head>
<body>
    <div><?= $bodyContent ?></div>
    <table>
        <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Desc</th>
        </tr>
    <?php foreach ($itemService->getItems() as $it): ?>
        <tr>
            <td><?= $it->id ?></td>
            <td><?= $it->name ?></td>
            <td><?= $it->desc ?></td>
        </tr>
    <?php endforeach ?>
    </table>
</body>
</html>
