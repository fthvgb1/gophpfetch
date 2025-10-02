<?php

use Fthvgb1\GoPHPFetch\Fetch;

include 'vendor/autoload.php';
include 'lib/profile.php';

$_ENV['caller'] = 'go';
echo 'gophpfetch:', PHP_EOL;
$a = Fetch::fetch(include 'requests.php');
print_r($a);
echo PHP_EOL;

