<?php

include 'vendor/autoload.php';
include 'lib/profile.php';
include 'lib/curl.php';

$_ENV['caller'] = 'swoole';
echo 'swoole curl:', PHP_EOL;
//Swoole\Runtime::enableCoroutine(SWOOLE_HOOK_CURL);
Co\run(function () {
    $curls = buildCurls(include 'requests.php');
    foreach ($curls as $ch => $request) {
        go(function () use ($ch, $request) {
            $t = curl_exec($ch);
            $request['callback']($t);
        });
    }
});
