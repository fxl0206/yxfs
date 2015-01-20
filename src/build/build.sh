IOUT=$1
cd ../
GOPATH=`pwd`
export GOPATH
cd $GOPATH
git pull
cp -R $GOPATH/src/views $IOUT
cp -R $GOPATH/src/conf $IOUT
go build -o $IOUT/yxfs $GOPATH/src/main.go
cp -R $GOPATH/src/build/start.sh $IOUT