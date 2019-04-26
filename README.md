# DNS over HTTPS
This client will make DNS over HTTPS request to a DNS server of your choosing. At this point it is only tested on MacOS.

## Config
You can set different endpoints in `config.json`, you also have to add a fingerprint of the public key of the endpoint.

## Run
- `./run.sh` or `nohup ./run.sh &`
- The client will ask for root privileges to start the DNS resolver on localhost.
- Set your DNS settings to 127.0.0.1
