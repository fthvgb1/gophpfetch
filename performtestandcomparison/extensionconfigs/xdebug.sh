#!/bin/bash
on="zend_extension=xdebug"
off=";$on"
if [ $1 != "on" ]; then
  s="$off"
  off="$on"
  on="$s"
  echo "s/$off/$on/"
fi
sed -i -e "s/$off/$on/" /etc/php/8.4/cli/conf.d/20-xdebug.ini
cat /etc/php/8.4/cli/conf.d/20-xdebug.ini