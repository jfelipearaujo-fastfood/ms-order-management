#!/bin/sh

localstack_url=http://localhost:4566
region=us-east-1
queue_name=UpdateOrderQueue

aws configure set profile localstack
aws configure set aws_access_key_id test
aws configure set aws_secret_access_key test
aws configure set region "$region"

queue_url=$(aws sqs get-queue-url --endpoint-url "$localstack_url" --output text --queue-name "$queue_name" --region "$region")

if [ $? -eq 0 ]; then
    echo "Queue URL: $queue_url"
    echo "Sending a message..."

    # message='{
    #     "order_id": "0674c12f-63ca-49e2-a2b2-d9113e97136e",
    #     "payment": {
    #         "id": "9f885861-a742-4e13-9eb2-c5d5e1a2627b",
    #         "state": "Approved"
    #     }
    # }'
    message='{
        "order_id": "0674c12f-63ca-49e2-a2b2-d9113e97136e",
        "order": {
            "state": "Received"
        }
    }'

    # Publish the message to the queue
    aws sqs send-message \
        --endpoint-url "$localstack_url" \
        --region "$region" \
        --queue-url "$queue_url" \
        --output text \
        --message-body "$message" > /dev/null

    # Check if the message publishing was successful
    if [ $? -eq 0 ]; then
        echo "Message published successfully."
    else
        echo "Failed to publish message."
    fi
else
    echo "Failed to retrieve the queue URL."
fi