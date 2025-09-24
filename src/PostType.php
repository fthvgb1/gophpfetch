<?php

declare(strict_types=1);

namespace Fthvgb1\GoPHPFetch;

enum PostType: string
{
    case FormUrlencoded = "x-www-form-urlencoded";
    case FormData = "form-data";
    case Json = "json";
    case Plain = 'plain';
    //case Binary = "binary";
}