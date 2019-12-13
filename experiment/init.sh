DIR=$(dirname $0)
IP_FILE=${DIR}/aws/ips.list

bash ${DIR}/utils/server_cmd.sh $IP_FILE 'mkdir -p ~/.tendermint/config'
bash ${DIR}/utils/server_cmd.sh $IP_FILE 'cp /home/ubuntu/go/src/github.com/tendermint/tendermint/experiment/data/genesis.json ~/.tendermint/config/genesis.json'
bash ${DIR}/utils/server_cmd.sh $IP_FILE 'tendermint init'
