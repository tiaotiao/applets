#!/bin/sh

count=$(ps -ax | grep 'centralserver' | grep -v 'grep' | wc -l)
if test $count -eq 0 
then
    echo 'server is not running'
else

    ps -ax | grep 'centralserver' | grep -v 'grep' | while read -r line
    do
        
        PID=$(echo $line | awk '{print $1}')
        name=${line#*config/}
        kill $PID
        echo server stoped $name
    done
fi

