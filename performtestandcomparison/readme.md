#### environment

- Linux 557bc0684e04 5.15.133.1-microsoft-standard-WSL2 #1 SMP Thu Oct 5 21:02:42 UTC 2023 x86_64 GNU/Linux

- php-cli 8.4
- xdebug 3.4.5
- xhprof 2.3.9

#### test

- edit extension config, set profile directory, reference [extensionconfigs](extensionconfigs)
- edit [profile.php](lib/profile.php) set [php-profiler](https://github.com/perftools/php-profiler) upload url
- edit [requests.php](requests.php) to your test request then execute

```shell
composer install && make