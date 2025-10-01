### a php http request extension on concurrence writing in golang

use FFI to call golang http requests with concurrence.

#### install

add follow contents to project's composer.json and run `composer install && composer run-script downloadExtension`.

```json
{
  "require": {
    "fthvgb1/gophpfetch": "*"
  },
  "scripts": {
    "downloadExtension": [
      "@putenv COMPOSER=vendor/fthvgb1/gophpfetch/composer.json",
      "@composer downloadExtension"
    ]
  }
}
```

to avoid including the unnecessary others platform's extension files, you need to download extension file manually
in [release](https://github.com/fthvgb1/gophpfetch/releases).

#### example

```php
$results = Fetch::fetch([
    [
        'url' => 'download url',
        'method' => 'get',
        'query' => ['query param' => 'value'],
        'id' => 'dowload file',
        'saveFile' => [
            'path' => 'save path',
            'mode' => '0644', //default
            'dirMode '=> '0755' //default
        ]
    ],
    [
        'url' => 'upload url',
        'method' => 'post',
        'id' => 'upload file',
        'header' => ['Content-Type' => PostType::FormData],
        'body' => [
            '__uploadFiles' => [
                'local file' => 'field'
            ]
        ]
    ]
]);

```

#### perform comparison

upload 5 files and download 5 files in a request on concurrence

results:

|        | gophpfetch                                               | multiple curl                                              | swoole curl                                              |
|--------|----------------------------------------------------------|------------------------------------------------------------|----------------------------------------------------------|
| xdebug | ![](performtestandcomparison\pics\xdebug-gophpfetch.png) | ![](performtestandcomparison\pics\xdebug-multiplecurl.png) | ![](performtestandcomparison\pics\xdebug-swoolecurl.png) |
| xhprof | ![](performtestandcomparison\pics\xhprof-overview.png)   |||
|        | ![](performtestandcomparison\pics\xhprof-gophpfetch.png) | ![](performtestandcomparison\pics\xhprof-multiplecurl.png) | ![](performtestandcomparison\pics\xhprof-swoolecurl.png) |

test method and scripts [see here](performtestandcomparison/readme.md)

#### reference

[dwdcth/php_golang_ffi](https://github.com/dwdcth/php_golang_ffi)