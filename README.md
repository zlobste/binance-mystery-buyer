# Mystery Box Buyer Bot

#### Bot is able to buy mystery boxes on Binance NFT platform.

## Configurations `config.yaml`

```
log: info # logging level
auth:
  proxy: ... # https://www.binance.com or proxy
  csrf_token: ... # CSRF token from your browser on www.binance.com
  cookie: ... # Cookie from your browser on www.binance.com
```

## Launch

#### Requires Docker or Golang

### Docker

```
# Install docker image and run the container
docker-compose pull && docker-compose up -d

# Display container logs
docker-compose logs -f
```

### Golang

```
# Set config file
export CONFIG=config.yaml

# Run the bot
binance-mystery-buyer run buyer
```
