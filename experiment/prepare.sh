DIR=$(dirname $0)
IP_FILE=${DIR}/aws/ips.list

bash ${DIR}/utils/server_cmd.sh $IP_FILE 'mkdir -p ~/go/src/github.com/tendermint'
bash ${DIR}/utils/server_cmd.sh $IP_FILE 'cd ~/go/src/github.com/tendermint && git clone https://zqzqz:fszqz123456@github.com/zqzqz/tendermint.git'
sleep 30
bash ${DIR}/utils/server_cmd.sh $IP_FILE 'bash ~/go/src/github.com/tendermint/tendermint/experiment/server/prepare.sh'
