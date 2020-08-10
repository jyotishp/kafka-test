# Tweets Service

## Usage

### Compile protobuf files
```bash
make proto
```

### Start the services
```bash
export AWS_REGION="us-east-1"
export DB_ENDPOINT="localhost:8000"
export DISABLE_SSL=true
make run
```

Swagger UI will be available at http://localhost:8080/swagger-ui

## What is there in the package?

### Tweets service
- A gRPC service to create add tweets
- This is proxied to `POST /v1/api/tweet` via grpc-gateway
- Pushes the tweets to a kafka stream

### Tweet Analyzer
- Reads the tweets from kafka stream
- Do a regular expression match on the message content
- If there is a match, add the tweet to DynamoDB
