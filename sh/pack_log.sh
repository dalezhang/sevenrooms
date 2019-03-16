#!/bin/bash
time=$(date -d '-1 day' +%Y-%m-%d)
cd /opt/sevenroom/log
cp  test.log packing.log
> test.log
find -name packing.log -exec tar -zcvf test.log.$time.tar.gz {} --remove-files \;

time2=$(date -d "7 days ago" +%Y-%m-%d)
find -name test.log.$time2.tar.gz -exec rm {} \;