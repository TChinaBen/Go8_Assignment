#!/bin/bash -ile

pid=$(ps -aux|grep m3b|grep -v 'grep' |awk '{print $2}')
if [[ $pid -ne "" ]];then
  kill -9 $pid
  echo "mimiclet stopped"
fi
