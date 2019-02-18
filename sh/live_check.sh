cd /opt/optitable_middleware
pid=$(cat optitable_middleware.pid)
if ps -ax | awk '{ print $1 }' | grep $pid
then
echo "Find process id : $pid. !!!!!"
fi
if ! ps -ax | awk '{ print $1 }' | grep $pid
then
echo "not find process id : $pid. !!!!!"
nohup RUN_MODE=staging ./main >> log/test.log 2>&1 & echo $! > optitable_middleware.pid
echo "boot main successfull $!"
fi