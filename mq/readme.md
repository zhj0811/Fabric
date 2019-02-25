## 初始化

```shell
#拉取docker镜像并启动
./start.sh
```

## 启动Rabbit容器

```shell
docker-compose up -d
```

## 构建模拟接收端(消费者)，并监听

```shell
#构建消费者
go build test_receive.go
#启动监听
./test_receive
```

