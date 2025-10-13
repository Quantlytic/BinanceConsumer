# GitHub Repository Secrets

This document outlines the GitHub repository secrets required for the BinanceConsumer deployment pipeline.

## Required Secrets

### Production Environment
- `KAFKA_BROKERS` - Kafka brokers connection string for production (e.g., "kafka-prod:9092")
- `KAFKA_TOPIC` - Kafka topic name for production (e.g., "binance-ticker")
- `KAFKA_CLIENT_ID` - Kafka client ID for production (e.g., "binance-consumer-prod")

### Development Environment
- `KAFKA_BROKERS_DEV` - Kafka brokers connection string for development (e.g., "kafka-dev:9092")
- `KAFKA_TOPIC_DEV` - Kafka topic name for development (e.g., "binance-ticker-dev")
- `KAFKA_CLIENT_ID_DEV` - Kafka client ID for development (e.g., "binance-consumer-dev")

## How to Set Secrets

1. Go to your GitHub repository
2. Navigate to Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Add each secret with the appropriate name and value

## Deployment Behavior

- **Production**: Deploys when pushing to `main` or `release/*` branches
- **Development**: Deploys when pushing to `dev` branch
- ConfigMaps are created dynamically from secrets during deployment
- The static `configmap.yaml` file has been removed and replaced with `configmap.yaml.example` for reference

## Security Notes

- Secrets are encrypted and only available to GitHub Actions
- Secrets are not displayed in logs
- Each environment uses separate secrets to maintain isolation