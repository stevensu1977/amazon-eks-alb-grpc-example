# eks-alb-grpc-example
This is example about how to use ALB ingress for grpc service in Amazon EKS cluster

### 前置条件

* 安装go, 笔者使用的是go1.15.6

* 安装protoc 工具(可选, 如果需要自己手动生成chat.pb.go)

* eksctl, kubectl 工具

* 已经创建了 Amazon EKS 集群,并且已经正确安装了VPC CNI, Amazon Loadbalancer Controller

  

### 1. 生成chat grpc golfing代码

```bash
#可选，chat.pb.go已经生成
protoc --go_out=plugins=grpc:chat chat/chat.proto
```

### 2. 编译服务器和客户端

```ba
#如果跨平台编译请配置GOOS=linux/drawin/windows

cd server 
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

cd client 
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
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

例子已经推送到stevensu/grpc-example-server,如果使用自己的docker image repo请修改yaml文件对应的id

```bash
docker build -t grpc-example-server .
docker tag grpc-example-server [docker hub id]/grpc-example-server
docker push [docker hub id]/grpc-example-server
```



### 5. 部署grpc server 到Amaon EKS集群

```bash
#这里需要注意一下，这里使用的是一个测试域名grpc.flowq.io, 请使用自己的域名,并且使用AWS ACM创建一个泛域名证书
kubectl apply -f grpc-example-server.yaml

#查看grpc server 是否正常工作
kubectl logs -f $(kubectl get pod | grep grpc-example-deployment | awk {'print $1'})

#参考输出
Go gRPC chat.ChatService start successful, listen on :50051

./code/client/client -host grpc.flowq.io -port 443 -tls

#参考输出,有Hello From this Server,则表示通过ALB Ingress转发到EKS的grpc-example-server 已经可以正常工作了
Endpoint: grpc.flowq.io:443,  TLS: true
2022/01/14 08:35:06 Response from server: Hello From the Server!
```





### 6. Tips

* ALB Ingress gRPC 健康检查代码,默认的是12, Amazon文档提到0~99, 找到另外一篇文章, https://grpc.github.io/grpc/core/md_doc_statuscodes.html, 世纪测试0,12, 或者0~16都可以
* Golang 客户端,初始化需要使用 grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}))) 初始化TLS才能正确调用ALB Ingress 
