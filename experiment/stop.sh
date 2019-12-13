DIR=$(dirname $0)
IP_FILE=${DIR}/aws/ips.list

source ${DIR}/config.sh
RESULT_DIR=${DIR}/results

bash ${DIR}/utils/server_cmd.sh $IP_FILE 'pkill /home/ubuntu/go/bin/tendermint'

IPS=$(cat $IP_FILE)

for IP in $IPS
do
    scp -i ${SSH_KEY} ${USERNAME}@${IP}:/home/${USERNAME}/.tendermint/run.log ${RESULT_DIR}/${IP}.log &
done
