# eks-alb-grpc-example
This is example about how to use ALB ingress for grpc service in Amazon EKS cluster

### Preconfig

* 安装go, 笔者使用的是go1.15.6
* 安装protoc 工具(可选, 如果需要自己手动生成chat.pb.go)

### 1. 生成chat grpc golfing代码

```bash
#可选，chat.pb.go已经生成
protoc --go_out=plugins=grpc:chat chat/chat.proto
```

### 2. 编译服务器和客户端

```ba
#如果跨平台编译请配置GOOS=linux/drawin/windows

cd server 
go build

cd client 
go build
```

* 需要注意的是使用ALB Ingress代理grpc服务只能走HTTPS/TLS协议， 所以client 启动的时候需要添加--tls 参数

```bash
  #go grpc client 需要使用TLS参数来初始化,否则会报错
  grpc.Dial(connectStr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
  
```

   

### 3. 本地测试

```bash
cd server 
./server 

输出:
Go gRPC chat.ChatService start successful, listen on :50051

#另外终端窗口
cd client 
./client

输出:
Endpoint: localhost:50051,  TLS: false
2022/01/14 15:28:30 Response from server: Hello From the Server!
```



### 4. 打包docker image

