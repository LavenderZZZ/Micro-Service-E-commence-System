#!/bin/bash

# 定义要创建的目录列表
DIRECTORIES=(
    "/mydata/mysql/data/db"
    "/mydata/mysql/data/conf"
    "/mydata/mysql/log"
    "/mydata/nginx/conf"
    "/mydata/nginx/html"
    "/mydata/nginx/log"
    "/mydata/rabbitmq/data"
    "/mydata/rabbitmq/log"
    "/mydata/elasticsearch/plugins"
    "/mydata/elasticsearch/data"
    "/mydata/logstash"
    "/mydata/mongo/db"
    "/mydata/app/mall-admin/logs"
    "/mydata/app/mall-portal/logs"
    "/mydata/app/mall-search/logs"
    "/mydata/app/mall-monitor/logs"
    "/mydata/app/mall-portal/logs"
    "/mydata/app/mall-auth/logs"
    "/mydata/app/mall-gateway/logs"
)

# 循环创建每一个目录
for DIR in "${DIRECTORIES[@]}"; do
    mkdir -p "$DIR"
    echo "Created directory: $DIR"
done

# 创建特定的文件
touch "/mydata/logstash/logstash.conf"
echo "Created file: /mydata/logstash/logstash.conf"

echo "All directories and files have been created."
