DIR=$(dirname $0)
IP_FILE=${DIR}/aws/ips.list

source ${DIR}/config.sh

bash ${DIR}/utils/server_cmd.sh $IP_FILE 'rm -rf ~/.tendermint'

IPS=$(cat $IP_FILE)
i=0

for IP in $IPS
do
	ssh -i ${SSH_KEY} ${USERNAME}@${IP} "cp -r /home/ubuntu/go/src/github.com/tendermint/tendermint/experiment/data/node${i} ~/.tendermint"
	i=$(expr "$i" + 1)
	ssh -i ${SSH_KEY} ${USERNAME}@${IP} "sed -e 's/192.198.1.${i}/${IP}/g' ~/.tendermint/config/config.toml"
done
