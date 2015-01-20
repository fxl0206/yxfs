for pid in `ps aux | grep "yxfs" | grep -v grep | awk '{print $2}'`
do
kill -9 $pid
echo "$pid killed success!"
cd ../
echo "yxfs conf/yx.ini  2>&1 | cronolog log/wss_%y%m%d.log &"
yxfs conf/yx.ini  2>&1 | cronolog log/wss_%y%m%d.log &
done