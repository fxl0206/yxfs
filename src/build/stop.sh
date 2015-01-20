for pid in `ps aux | grep "yxfs" | grep -v grep | awk '{print $2}'`
do
kill -9 $pid
echo "$pid killed success!"
done
