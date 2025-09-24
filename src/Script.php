<?php

declare(strict_types=1);

namespace Fthvgb1\GoPHPFetch;

use Exception;
use RuntimeException;

class Script
{
    const VERSION = '0.2.2';

    /**
     * @throws Exception
     */
    public static function checkAndDownloadExtension(): void
    {
        $file = self::getPlatformArchExtension();
        $version = self::VERSION;
        if (!file_exists($file)) {
            $filename = basename($file);
            $url = "https://github.com/fthvgb1/gophpfetch/releases/download/v{$version}/$filename";
            echo 'start to download ', $url, PHP_EOL;
            $data = file_get_contents($url);
            if (false === $data) {
                throw new Exception("can't download extension file: {$url}");
            }
            $r = file_put_contents($file, $data);
            if (false === $r) {
                throw new Exception("write extension {$file} failed!");
            }
            echo 'completed downloading extension', PHP_EOL;
        }
    }

    /**
     * @return string
     */
    public static function getPlatformArchExtension(): string
    {
        $os = strtolower(PHP_OS_FAMILY);
        $archMap = [
            'linux' => 'linux_x86_64.so',
            'windows' => 'windows_x86_64.dll',
            'darwin' => 'darwin_' . php_uname('m') . '.dylib'
        ];
        if (!isset($archMap[$os])) {
            throw new RuntimeException('not support this platform or arch');
        }
        $filename = 'gophpfetch_' . Script::VERSION . '_' . $archMap[$os];
        return __DIR__ . '/exts/' . $filename;
    }
}