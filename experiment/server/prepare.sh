set -e

DIR=$(dirname $0)
source ${DIR}/../config.sh

sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt-get update
sudo apt-get install -y golang-go wget git build-essential software-properties-common curl vim

echo "\nGOROOT=/usr/local/go\n" >> ~/.bashrc
echo "\nGOPATH=/home/ubuntu/go/\n" >> ~/.bashrc
echo '\nPATH=$GOPATH/bin:$GOROOT/bin:$PATH\n' >> ~/.bashrc

cd $WORKSPACE && make install
