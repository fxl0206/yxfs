IOUT=$1
cd ../../
GOPATH=`pwd`
export GOPATH
cd $GOPATH
git pull
cp -R $GOPATH/src/html $IOUT
cp -R $GOPATH/src/conf $IOUT
go build -o $IOUT/yxfs $GOPATH/src/main.go
chmod +x $GOPATH/src/build/*.sh
cp -R $GOPATH/src/build/*.sh $IOUT