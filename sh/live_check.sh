cd /opt/sevenroom
pid=$(cat sevenrooms.pid)
if ps -ax | awk '{ print $1 }' | grep $pid
then
echo "Find process id : $pid. !!!!!"
fi
if ! ps -ax | awk '{ print $1 }' | grep $pid
then
echo "not find process id : $pid. !!!!!"
RUN_MODE=staging nohup ./main >> log/test.log 2>&1 & echo $! > sevenrooms.pid
echo "boot main successfull $!"
fi