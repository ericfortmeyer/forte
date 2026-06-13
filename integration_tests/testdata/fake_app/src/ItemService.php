<?php

class ItemService
{
    public function getItems(): array
    {
        return [
            (object) [
                "id" => "1",
                "name" => "it",
                "desc" => "does things"
                ],
                (object) [
                "id" => "2",
                "name" => "that",
                "desc" => "does other things"
                ],
                (object) [
                "id" => "3",
                "name" => "other thing",
                "desc" => "doesn't do much"
            ],
        ];
    }
}
