while true
do
	random=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
	curl -s "localhost:26657/broadcast_tx_async?tx=\"${random}\"" > /dev/null
done
