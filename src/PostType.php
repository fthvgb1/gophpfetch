<?php

declare(strict_types=1);

namespace GoPHP\GoPHPFetch;

enum PostType: string
{
    case FormUrlencoded = "x-www-form-urlencoded";
    case FormData = "form-data";
    case Json = "json";
    //case Binary = "binary";
}