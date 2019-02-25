#!/usr/bin/env bash
docker-compose up -d


docker exec rabbitmq_rabbit1_1 rabbitmqctl set_policy ha-all "" '{"ha-mode":"all","ha-sync-mode":"automatic"}'
