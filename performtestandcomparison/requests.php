<?php

use Fthvgb1\GoPHPFetch\PostType;

return [
    [
        'url' => 'http://192.168.43.229:17778/upload',
        'id' => 'upload-1',
        'method' => 'post',
        'header' => ['Content-Type' => PostType::FormData],
        'body' => [
            '__uploadFiles' => [
                'imgs/1.jpg' => 'uploads/1-' . $_ENV['caller'] . '-05.jpg'
            ]
        ]
    ],
    [
        'url' => 'http://192.168.43.229:17778/upload',
        'method' => 'post',
        'id' => 'upload-2',
        'header' => ['Content-Type' => PostType::FormData],
        'body' => [
            '__uploadFiles' => [
                'imgs/2.jpg' => 'uploads/2-' . $_ENV['caller'] . '.jpg'
            ]
        ]
    ],
    [
        'url' => 'http://192.168.43.229:17778/upload',
        'method' => 'post',
        'id' => 'upload-3',
        'header' => ['Content-Type' => PostType::FormData],
        'body' => [
            '__uploadFiles' => [
                'imgs/3.jpg' => 'uploads/3-' . $_ENV['caller'] . '.jpg'
            ]
        ]
    ],
    [
        'url' => 'http://192.168.43.229:17778/upload',
        'id' => 'upload-4',
        'method' => 'post',
        'header' => ['Content-Type' => PostType::FormData],
        'body' => [
            '__uploadFiles' => [
                'imgs/4.jpg' => 'uploads/4-' . $_ENV['caller'] . '.jpg'
            ]
        ]
    ],
    [
        'url' => 'http://192.168.43.229:17778/upload',
        'id' => 'upload-5',
        'method' => 'post',
        'header' => ['Content-Type' => PostType::FormData],
        'body' => [
            '__uploadFiles' => [
                'imgs/5.jpg' => 'uploads/5-' . $_ENV['caller'] . '.jpg'
            ]
        ]
    ],
    [
        'url' => 'http://192.168.43.229:12345/db7d189a0370afc08bd18.jpg',
        'id' => 'download-1',
        'saveFile' => [
            'path' => $_ENV['caller'] . '/db7d189a0370afc08bd18.jpg'
        ]
    ],
    [
        'url' => 'http://192.168.43.229:12345/2c69f642f4cdf2f7015f0.jpg',
        'id' => 'download-2',
        'saveFile' => [
            'path' => $_ENV['caller'] . '/2c69f642f4cdf2f7015f0.jpg'
        ]
    ],
    [
        'url' => 'http://192.168.43.229:12345/95758908cb18c032f76ff.jpg',
        'id' => 'download-3',
        'saveFile' => [
            'path' => $_ENV['caller'] . '/95758908cb18c032f76ff.jpg'
        ]
    ],
    [
        'url' => 'http://192.168.43.229:12345/ef898dced4dc5dc06312e.jpg',
        'id' => 'download-4',
        'saveFile' => [
            'path' => $_ENV['caller'] . '/ef898dced4dc5dc06312e.jpg'
        ]
    ],
    [
        'url' => 'http://192.168.43.229:12345/6373811d966ca9d171a0e.jpg',
        'id' => 'download-5',
        'saveFile' => [
            'path' => $_ENV['caller'] . '/6373811d966ca9d171a0e.jpg'
        ]
    ],
];