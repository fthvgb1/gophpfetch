<?php


use FFI\CData;

class fetch
{
    private static self $instance;
    public FFI $ffi;

    public function __construct()
    {
        $this->ffi = FFI::load('gophpfetch.h');
        self::$instance = $this;
    }


    public static function init(): void
    {
        self::$instance = new self();
    }

    /**
     * @param array{
     *     requests: array{
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
     * }[]
     * ,concurrency:int} $arr timeout Millisecond
     * @return array{
     *     res:array{requestId:string,
     *     header:array<string,string>,
     *     httpStatusCode:int,
     *     res:string
     *     }[], err:string}
     */
    public static function fetch(array $arr): array
    {
        $confStr = json_encode($arr);
        $conf = self::$instance->makeGoString($confStr);
        $r = self::readCString(self::$instance->ffi->Fetch($conf));
        return json_decode($r, true);
    }

    public function makeGoString(string $str): CData
    {
        $goStr = self::$instance->ffi->new('GoString', 0);
        $size = strlen($str);
        $cStr = FFI::new("char[$size]", 0);
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