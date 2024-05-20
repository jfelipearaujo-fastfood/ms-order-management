#!/bin/sh

echo "Initializing Secrets Manager..."

awslocal secretsmanager create-secret \
    --name db-secret-url \
    --description "DB URL" \
    --secret-string "postgres://order:order@test:5432/order_db?sslmode=disable"