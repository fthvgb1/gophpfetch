<?php
include 'fetch.php';


/**
 * @var $a array{err:string,res:array{}[]}
 */
$a = fetch::fetch([
    'requests' => [
        [
            'url' => 'http://192.168.43.229:8765/',
            'method' => 'get',
            'timeout' => 5000
        ],
        [
            'url' => '',
        ]
    ],
]);
print_r($a);