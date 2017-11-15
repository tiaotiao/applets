#!/bin/sh

isrunning()
{
    name=$1
    count=$(ps -ax | grep 'centralserver' | grep $name | grep -v 'grep' | wc -l)
    
    if test $count -eq 0
    then
        return 0
    else
        return 1
    fi
}

for file in ./config/*
do
    if test -f $file
    then
        config=$(echo $file | grep 'server-')
        if [[ $config == "" ]] 
        then
            continue
        fi

        name=${config#*config/}

        isrunning $name

        running=$?
        if test $running -eq 1
        then
            echo server is running $name
            continue
        fi

        nohup ./bin/centralserver -config=$config >/dev/null 2>&1 &
        echo start server OK $name
    fi
done