<?php

declare(strict_types=1);

namespace GoPHP\GoPHPFetch;

use FFI;
use FFI\CData;


class Fetch
{

    private static self $instance;

    /**
     * @return Fetch
     */
    public static function getInstance(): Fetch
    {
        return self::$instance;
    }

    public FFI $ffi;
    public array $freeVars;

    public function __construct()
    {
        $this->ffi = FFI::cdef(<<<CT
typedef long int ptrdiff_t;
typedef long unsigned int size_t;
typedef int wchar_t;
typedef struct {
  long long __max_align_ll __attribute__((__aligned__(__alignof__(long long))));
  long double __max_align_ld __attribute__((__aligned__(__alignof__(long double))));
} max_align_t;
typedef struct { const char *p; ptrdiff_t n; } _GoString_;
typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef size_t GoUintptr;
typedef float GoFloat32;
typedef double GoFloat64;
typedef char _check_for_64_bit_pointer_matching_GoInt[sizeof(void*)==64/8 ? 1:-1];
typedef _GoString_ GoString;
typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;
extern char* Fetch(GoString s, GoInt concurrence, GoUint8 associate);
CT, Script::getPlatformArchExtension());
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

Fetch::init();