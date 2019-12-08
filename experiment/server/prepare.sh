sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt-get update
sudo apt-get install golang-go

echo "\nGOROOT=/usr/local/go\n" >> ~/.bashrc
echo "\nGOPATH=/home/ubuntu/go/\n" >> ~/.bashrc
echo '\nPATH=$GOPATH/bin:$GOROOT/bin:$PATH\n' >> ~/.bashrc

cd ~/go/src/github.com/tendermint/tendermint && make install
