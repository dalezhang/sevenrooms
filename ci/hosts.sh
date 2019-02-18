#!/bash
echo "" >> /etc/hosts
echo "$HOST_IP mysql" >> /etc/hosts
cat /etc/hosts
