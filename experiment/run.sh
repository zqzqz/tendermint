DIR=$(dirname $0)
IP_FILE=${DIR}/aws/ips.list

bash ${DIR}/utils/server_cmd.sh $IP_FILE 'bash ~/go/src/github.com/tendermint/tendermint/experiment/server/run.sh'
