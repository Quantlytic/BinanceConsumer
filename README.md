# BinanceConsumer
Consumes live websocket feeds from binance and writes to kafka topic.

## Configuration

This application uses environment variables for configuration. See [`docs/SECRETS.md`](docs/SECRETS.md) for information about required GitHub repository secrets for deployment.

### Environment Variables

- `KAFKA_BROKERS` - Kafka brokers connection string (default: "localhost:9092")
- `KAFKA_TOPIC` - Kafka topic name (default: "binance-ticker")  
- `KAFKA_CLIENT_ID` - Kafka client ID (default: "binance-publisher")

## Deployment

The application is automatically deployed using GitHub Actions:
- **Development**: Push to `dev` branch
- **Production**: Push to `main` or `release/*` branches

Configuration is managed through GitHub repository secrets. See [`docs/SECRETS.md`](docs/SECRETS.md) for setup instructions.
