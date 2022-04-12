#!/bin/bash -ile

dir="$(cd "$(dirname "$0")" && pwd)"

cd $dir
m3b=$(ls |grep m3b)

chmod +x $m3b

if [[ $(pidof $m3b) != "" ]];then
  kill -9 $(pidof $m3b)
  echo "restart mimiclet"
fi

nohup ./$m3b >>log.txt 2>&1 &
echo "mimiclet started: $m3b"
