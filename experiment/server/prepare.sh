set -e

DIR=$(dirname $0)
source ${DIR}/../config.sh

sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt-get update
sudo apt-get install -y golang-go wget git build-essential software-properties-common curl vim

echo "GOROOT=/usr/local/go" >> ~/.bashrc
echo "GOPATH=/home/ubuntu/go/" >> ~/.bashrc
echo 'PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.bashrc

cd $WORKSPACE && make install
