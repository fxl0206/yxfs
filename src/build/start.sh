sh stop.sh
echo "yxfs conf/yx.ini  2>&1 | cronolog log/wss_%y%m%d.log &"
yxfs conf/yx.ini  2>&1 | cronolog log/wss_%y%m%d.log &
