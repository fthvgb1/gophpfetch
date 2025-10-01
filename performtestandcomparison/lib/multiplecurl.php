<?php

include 'curl.php';

/**
 * @throws Exception
 */
function multipleCurl(array $requests): void
{
    $mh = curl_multi_init();
    $curls = buildCurls($requests, true);
    foreach ($curls as $ch => $req) {
        curl_multi_add_handle($mh, $ch);
    }

    do {
        $status = curl_multi_exec($mh, $unfinishedHandles);
        //curl_multi_select($mh);
        tackle($mh, $status, $curls, $unfinishedHandles);

    } while ($unfinishedHandles);

    curl_multi_close($mh);
}


/**
 * @throws Exception
 */
function tackle($mh, $status, $curls, $unfinishedHandles): void
{
    if ($status !== CURLM_OK) {
        throw new \Exception(curl_multi_strerror(curl_multi_errno($mh)));
    }

    while (($info = curl_multi_info_read($mh)) != false) {
        if ($info['msg'] === CURLMSG_DONE) {
            $handle = $info['handle'];
            curl_multi_remove_handle($mh, $handle);
            $curls[$handle]['callback']($info);
        }
    }
    if ($unfinishedHandles && curl_multi_select($mh) === -1) {
        throw new \Exception(curl_multi_strerror(curl_multi_errno($mh)));
    }
}


/**
 * @throws Exception
 */
function multipleCurlByFiber($requests): void
{
    $curlHandles = [];
    $curls = buildCurls($requests, true);
    $mh = curl_multi_init();
    $mh_fiber = curl_multi_init();
    $halfOfList = floor(count($curls) / 2);
    $index = 0;
    foreach ($curls as $ch => $request) {
        $curlHandles[] = $ch;
        // half of urls will be run in background in fiber
        $index > $halfOfList ? curl_multi_add_handle($mh_fiber, $ch) : curl_multi_add_handle($mh, $ch);
        $index++;
    }
    unset($index);

    $fiber = new Fiber(function (CurlMultiHandle $mh) use ($curls) {
        $still_running = null;
        do {
            $status = curl_multi_exec($mh, $still_running);
            tackle($mh, $status, $curls, $still_running);
            Fiber::suspend();
        } while ($still_running);
    });

// run curl multi exec in background while fiber is in suspend status
    $fiber->start($mh_fiber);

    $still_running = null;
    do {
        $status = curl_multi_exec($mh, $still_running);
        tackle($mh, $status, $curls, $still_running);
    } while ($still_running);

    do {
        /**
         * at this moment curl in fiber already finished (maybe)
         * so we must refresh $still_running variable with one more cycle "do while" in fiber
         **/
        $status_fiber = $fiber->resume();
    } while (!$fiber->isTerminated());

    foreach ($curlHandles as $index => $ch) {
        $index > $halfOfList ? curl_multi_remove_handle($mh_fiber, $ch) : curl_multi_remove_handle($mh, $ch);
    }
    curl_multi_close($mh);
    curl_multi_close($mh_fiber);
}

