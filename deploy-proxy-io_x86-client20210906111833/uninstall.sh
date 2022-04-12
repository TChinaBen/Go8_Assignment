#!/bin/bash -il

CLIENT_PATH="$(cd "$(dirname "$0")" && pwd)"
echo $CLIENT_PATH

SVC_PATH=/etc/systemd/system/mimiclet.service

if [ -f $SVC_PATH ];then
    systemctl disable mimiclet
    systemctl stop mimiclet
    systemctl status mimiclet
    if [ $? -eq 0 ];then
        echo "mimiclet stop failed!"
    else
        echo "mimiclet stopped!"
    fi
    rm -f $SVC_PATH
    systemctl daemon-reload
    systemctl reset-failed
fi

cd $CLIENT_PATH
sh stop.sh
echo "mimiclet uninstalled!"
