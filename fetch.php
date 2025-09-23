<?php


use FFI\CData;

enum PostType: string
{
    case FormUrlencoded = "x-www-form-urlencoded";
    case FormData = "form-data";
    case Json = "json";
    //case Binary = "binary";
}

class fetch
{

    private static self $instance;

    /**
     * @return fetch
     */
    public static function getInstance(): fetch
    {
        return self::$instance;
    }

    public FFI $ffi;
    public array $freeVars;

    public function __construct()
    {
        $this->ffi = FFI::load('gophpfetch.h');
        self::$instance = $this;
    }

    public function free(): void
    {
        array_walk($this->freeVars, fn($r) => FFI::free($r));
        $this->freeVars = [];
    }

    public static function init(): void
    {
        self::$instance = self::$instance ?? new self();
    }

    /**
     * @param array<array{
     *     id:string,
     *     url:string,
     *     Method:string,
     *     query:array<string,mixed>,
     *     header: array<string,string>,
     *     body: array<string,string>,
     *     maxRedirectNum:int,
     *     timeout:int,
     *     saveFilename:string,
     *     getResponseHeader:bool
     * }> $arr timeout unit is Millisecond
     * contentType PostType
     * @param int $concurrence
     * @param bool $associate
     *
     * @return array{
     *     results: array<array{requestId:string,
     *     header:array<string,string>,
     *     httpStatusCode:int,
     *     result:string
     *     }> | array<string, array{requestId:string,
     *     header:array<string,string>,
     *     httpStatusCode:int,
     *     result:string
     *     }>, err:string}
     */

    public static function fetch(array $arr, int $concurrence = 0, bool $associate = false): array
    {
        foreach ($arr as &$item) {
            if (!isset($item['header']['Content-Type'])) {
                continue;
            }
            if ($item['header']['Content-Type'] !== PostType::Json) {
                continue;
            }
            if (isset($item['body']['jsonData'])) {
                continue;
            }
            $item['body']['jsonData'] = json_encode($item['body']);
        }
        unset($item);
        $requests = json_encode($arr);
        $conf = self::makeGoString($requests);
        $r = self::readCString(self::$instance->ffi->Fetch($conf, $concurrence, $associate));
        return json_decode($r, true);
    }

    public static function makeGoString(string $str): CData
    {
        $goStr = self::$instance->ffi->new('GoString', false);
        $size = mb_strlen($str);
        $cStr = self::$instance->ffi->new("char[$size]", false);
        FFI::memcpy($cStr, $str, $size);
        $goStr->p = $cStr;
        $goStr->n = strlen($str);
        self::$instance->freeVars[] = $goStr;
        self::$instance->freeVars[] = $cStr;
        return $goStr;
    }

    public static function readCString(CData $cData): string
    {
        $r = FFI::string($cData);
        FFI::free($cData);
        self::$instance->free();
        return $r;
    }
}

fetch::init();