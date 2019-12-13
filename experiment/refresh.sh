DIR=$(dirname $0)
IP_FILE=${DIR}/aws/ips.list

bash ${DIR}/utils/server_cmd.sh $IP_FILE 'cd ~/go/src/github.com/tendermint/tendermint && git reset --hard && git checkout dagchain && git pull origin dagchain'
