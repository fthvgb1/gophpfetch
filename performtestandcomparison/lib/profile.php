<?php

use Xhgui\Profiler\Profiler;

$profiler = new Xhgui\Profiler\Profiler([
    'save.handler' => Profiler::SAVER_STACK,
    'save.handler.stack' => array(
        'savers' => array(
            Profiler::SAVER_UPLOAD,
            Profiler::SAVER_FILE,
        ),
        // if saveAll=false, break the chain on successful save
        // 'saveAll' => false,
    ),
    'save.handler.file' => array(
        'filename' => '/tmp/' . $_SERVER["SCRIPT_FILENAME"] . '-xhgui.data.jsonl',
    ),
    'save.handler.upload' => array(
        'url' => 'http://192.168.43.229:13333/run/import',
        // The timeout option is in seconds and defaults to 3 if unspecified.
        'timeout' => 3,
        // the token must match 'upload.token' config in XHGui
        'token' => 'xhgui',
        // verify option to disable ssl verification, defaults to true if unspecified.
        'verify' => true,
    ),
]);

$profiler->start();
