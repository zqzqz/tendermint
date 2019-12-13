DIR=$(dirname $0)
IP_FILE=${DIR}/aws/ips.list

source ${DIR}/config.sh

bash ${DIR}/utils/server_cmd.sh $IP_FILE 'rm -rf ~/.tendermint'

IPS=$(cat $IP_FILE)
i=0
for IP in $IPS
do
	ssh -i ${SSH_KEY} ${USERNAME}@${IP} "cp -r /home/ubuntu/go/src/github.com/tendermint/tendermint/experiment/data/node${i} ~/.tendermint"
	i2=1
	for IP2 in $IPS
	do
		ssh -i ${SSH_KEY} ${USERNAME}@${IP} "sed -i 's/192.168.1.${i2}/${IP2}/g' ~/.tendermint/config/config.toml"
		i2=$(expr "$i2" + 1)
	done
	i=$(expr "$i" + 1)
done
