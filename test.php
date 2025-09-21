<?php

include 'fetch.php';


$a = fetch::fetch([
    [
        'url' => 'http://192.168.43.229:8765/',
        'method' => 'get',
        'timeout' => 5000
    ],
    [
        'url' => '',
    ]
]);
print_r($a);