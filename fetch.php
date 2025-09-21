<?php


use FFI\CData;

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

    public function __construct()
    {
        $this->ffi = FFI::load('gophpfetch.h');
        self::$instance = $this;
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
        $requests = json_encode($arr);
        $conf = self::$instance->makeGoString($requests);
        $r = self::readCString(self::$instance->ffi->Fetch($conf, $concurrence, $associate));
        return json_decode($r, true);
    }

    public function makeGoString(string $str): CData
    {
        $goStr = self::$instance->ffi->new('GoString', false);
        $size = strlen($str);
        $cStr = FFI::new("char[$size]", false);
        FFI::memcpy($cStr, $str, $size);
        $goStr->p = $cStr;
        $goStr->n = strlen($str);
        return $goStr;
    }

    public static function readCString(CData $cData): string
    {
        $r = FFI::string($cData);
        FFI::free($cData);
        return $r;
    }
}

fetch::init();