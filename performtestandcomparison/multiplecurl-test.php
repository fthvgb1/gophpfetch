<?php

include_once 'vendor/autoload.php';
include 'lib/profile.php';
include 'lib/multiplecurl.php';

$_ENV['caller'] = 'php';
echo 'multiple curl:', PHP_EOL;
multipleCurl(include 'requests.php');
//curlByFiber(include 'requests.php');

