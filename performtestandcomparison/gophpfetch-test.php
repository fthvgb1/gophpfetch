<?php

use Fthvgb1\GoPHPFetch\Fetch;

include 'vendor/autoload.php';
include 'lib/profile.php';

$_ENV['caller'] = 'go';
$a = Fetch::fetch(include 'requests.php', 0, 0);
print_r($a);

