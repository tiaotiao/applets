#!/bin/sh

count=$(ps -ax | grep 'centralserver' | grep -v 'grep' | wc -l)
if test $count -eq 0 
then 
    echo 'server is not running'
else
    ps -ax | grep 'centralserver' | grep -v 'grep'
fi
