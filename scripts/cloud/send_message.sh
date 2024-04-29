#!/bin/sh

localstack_url=http://localhost:4566
queue_name=UpdateOrderQueue

export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test

queue_url=$(aws sqs get-queue-url --endpoint-url "$localstack_url" --output text --queue-name "$queue_name")

if [ $? -eq 0 ]; then
    echo "Queue URL: $queue_url"
    echo "Sending a message..."

    # state `2` means the order was payed and is ready to be started by the kitchen
    message='{
        "order_id": "679e0f4d-799d-4bc1-8820-c49ed799ad6c",
        "state": 2
    }'

    # Publish the message to the queue
    aws sqs send-message \
        --endpoint-url "$localstack_url" \
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