<?php

use Fthvgb1\GoPHPFetch\PostType;

function buildCurls(array $requests, $multiple = false): WeakMap
{

    $curls = new WeakMap();
    foreach ($requests as $request) {
        $ch = curl_init($request['url']);
        $opts = [
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_MAXREDIRS => 10,
            CURLOPT_TIMEOUT => 0,
            CURLOPT_FOLLOWLOCATION => true,
            CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
        ];
        $opts[CURLOPT_CUSTOMREQUEST] = strtoupper($request['method'] ?? 'get');
        if (isset($request['header']['Content-Type']) && $request['header']['Content-Type']) {
            switch ($request['header']['Content-Type']) {
                case PostType::Json:
                    $opts[CURLOPT_HTTPHEADER] = ['Content-Type: application/json'];
                    $opts[CURLOPT_POSTFIELDS] = $request['body']['__Data'] ?? json_encode($request['body']);
                    break;
                case PostType::Plain:
                    $opts[CURLOPT_HTTPHEADER] = ['Content-Type: text/plain'];
                    $opts[CURLOPT_POSTFIELDS] = $request['body']['__Data'];
                    break;
                case PostType::FormData:
                    foreach ($request['body']['__uploadFiles'] as $local => $target) {
                        $opts[CURLOPT_POSTFIELDS][$target] = new CURLFile($local);
                    }
                    unset($request['body']['__uploadFiles']);
                    $opts[CURLOPT_POSTFIELDS] = array_merge($request['body'], $opts[CURLOPT_POSTFIELDS]);
                    break;
                case PostType::FormUrlencoded:
                    $opts[CURLOPT_HTTPHEADER] = "application/x-www-form-urlencoded";
                    $opts[CURLOPT_POSTFIELDS] = $request['body'];
                    break;
            }

        }

        curl_setopt_array($ch, $opts);
        if (!isset($request['callback'])) {
            if ($multiple) {
                $request['callback'] = function ($info) use ($ch, $request) {
                    $id = $request['id'] ?? $request['url'];
                    if ($info['result'] !== CURLE_OK) {
                        echo $id, ' happened error:', curl_error($ch), ' ', curl_strerror($info['result']), PHP_EOL;
                        return;
                    }
                    $statusCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
                    if (isset($request['saveFile']) && $request['saveFile']) {
                        $dir = dirname($request['saveFile']['path']);
                        if (!is_dir($dir)) {
                            mkdir($dir, 0755, true);
                        }
                        file_put_contents($request['saveFile']['path'], curl_multi_getcontent($ch));
                        echo $request['id'] ?? $request['url'], ' completed', PHP_EOL;
                        return;
                    }
                    $t = curl_multi_getcontent($ch);
                    echo $id, '=>', $statusCode, ' ', $t, PHP_EOL;
                };
            } else {
                $request['callback'] = function ($res) use ($ch, $request) {
                    $id = $request['id'] ?? $request['url'];
                    if (isset($request['saveFile']) && $request['saveFile']) {
                        $dir = dirname($request['saveFile']['path']);
                        if (!is_dir($dir)) {
                            mkdir($dir, 0755);
                        }
                        file_put_contents($request['saveFile']['path'], $res);
                        echo $request['id'] ?? $request['url'], ' completed', PHP_EOL;
                        return;
                    }
                    $statusCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
                    echo $id, '=>', $statusCode, ' ', $res, PHP_EOL;
                };
            }
        }

        $curls[$ch] = $request;
    }
    return $curls;
}