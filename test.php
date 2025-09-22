<?php

include 'fetch.php';


$a = fetch::fetch([
    [
        'url' => 'http://192.168.43.229:8765/',
        'method' => 'get',
        'timeout' => 5000
    ],
    [
        'url' => 'http://192.168.43.229:8765',
        'id' => 'query',
        'method' => 'post',
        'header' => ['Content-Type' => PostType::Json],
        'body' => [
            "action" => "findNotes",
            "version" => 6,
            "params" => [
                "query" => 'deck:生词 "正面:proactive inactive proactively inactively interactive interactively interactivity"'
            ]
        ]
    ],
    [
        'url' => 'http://192.168.43.229:12333/server.php',
        'id' => 'upload',
        'method' => 'post',
        'header' => ['Content-Type' => PostType::FormData],
        'body' => [
            'upload' => 'hello php',
            '__uploadFiles' => [
                './pic.jpg' => 'aa.jpg'
            ]
        ]
    ],
    [
        'url' => 'http://192.168.43.229:12333/server.php',
        'method' => 'get',
        'query' => ['download' => 'uploads/pic.jpg'],
        'id' => 'saveFile',
        'saveFilename' => 'uploads/bb.jpg|0644'
    ]
], 0, true);
print_r($a);