nohup /home/ubuntu/go/bin/tendermint node --proxy_app=kvstore --trace 2>&1 > ~/.tendermint/run.log &
echo "started"
