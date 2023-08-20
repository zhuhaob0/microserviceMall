## 依赖组件
- 基础组件
    - redis
    - zookeeper
    - consul
- 可视化组件
    - zipkin
    - hystrix-dashboard

## 项目启动
- 启动consul：`consul agent -dev`
- 启动zookeeper容器：`docker run -d -p 2181:2181 --name zookeeper zookeeper`
- 启动zipkin容器：`docker run -d -p 9411:9411 openzipkin/zipkin`
- 启动文件配置服务器：`go build && ./configServer`
- 启动sk-admin管理员模块：`go build && ./sk-admin`
- 启动网关gateway模块：`go build && ./gateway`
- 启动user-service模块：`go build && ./user-service`
- 启动oauth-service模块：`go build && ./oauth-service`
- 启动sk-app秒杀业务模块：`go build && ./sk-app`
- 启动sk-core秒杀内核模块：`go build && ./sk-core`


## 毕业设计API文档
- 配置文件服务器
    - api: 
        - 获取文件（GET）: `127.0.0.1:10085/master/sk-admin-dev.yaml`
    （注：路径的最后一部分是是文件名）
        - return: 文件内容的base64编码
        - 健康检查（GET）：`127.0.0.1:10085/health`
        - metrics: `127.0.0.1:10085/metrics`
- sk-admin管理员模块
    - api:
        - 创建商品（POST）：`127.0.0.1:9030/product/create`
        - 列出商品（GET）：`127.0.0.1:9030/product/list`
        - 修改商品（POST）：`127.0.0.1:9030/product/update`
        - 删除商品（POST）：`127.0.0.1:9030/product/delete`
        - 创建活动（POST）：`127.0.0.1:9030/activity/create`
        - 列出活动（GET）：`127.0.0.1:9030/activity/list`
        - 健康检查（GET）：`127.0.0.1:9030/health`
        - metrics: `127.0.0.1:9030/metrics`

- user-service 用户服务模块
    - api:
        - 创建用户（POST）：`127.0.0.1:9009/user/create`
        - 检查用户（POST）：`127.0.0.1:9009/user/check`
        - 健康检查（GET）：`127.0.0.1:9009/health`
        - metrics：`127.0.0.1:9009/metrics`

- gateway 网关
    - api:
        - 反向代理：`127.0.0.1:9090/sk-admin/product/list`
        （sk-admin是服务名，后边是服务请求，代理通过服务名在consul中寻找服务实例，获取它的服务地址，进行请求转发）

- oauth-service鉴权模块
    - api：

- sk-app模块
    - api：
        - 查询所有合法活动（GET）：`127.0.0.1:9031/sec/list`
        - 根据product_id查询活动（POST）：`127.0.0.1:9031/sec/info`
        - 商品秒杀（POST）：`127.0.0.1:9031/sec/kill`



## Postman调试记录
- getProductList
    - api: `127.0.0.1:9030/product/list`
    - method: `GET`

- create product
    - api: `127.0.0.1:9030/product/create`
    - method: `POST`
    - 在body的raw里边写请求内容
        ```json
        {
            "product_id": 1,
            "product_name": "华为P30",
            "total": 100,
            "status": 1
        }
        ```

- update product
    - api: `127.0.0.1:9030/product/update`
    - method: `POST`
    - 在body的raw里边写请求内容
        ```json
        {
            "product_name": "红米K40",
            "total": 200,
            "status": 2
        }
        ```

- delete product
    - api: `127.0.0.1:9030/product/delete`
    - method: `POST`
    - 在body的raw里边写请求内容
        ```json
        {
            "product_name": "红米K40",
            "total": 100,
            "status": 1
        }
        ```
---
- getActivityList
    - api: `127.0.0.1:9030/activity/list`
    - method: `GET`

- create activity
    - api: `127.0.0.1:9030/activity/create`
    - method: `POST`
    - 在body的raw里边写请求内容
        ```json
        {
            "activity_id": 1,
            "activity_name":"华为P30抢购日",
            "product_id":1,
            "start_time":1000,
            "end_time":10000,
            "total":100,
            "status":0,
            "speed":10,
            "buy_limit":30,
            "buy_rate":0.3
        }
        ```
- update activity

- delete activity

---
- create-user
    - api: `127.0.0.1:9009/user/create`
    - method: `POST`
    - 在body里的raw里边写内容
        ```json
        {
            "user_name":"aoho",
            "password":"123456",
            "age":22    
        }
        ```
- check-user
    - api: `127.0.0.1:9009/user/check`
    - method: `POST`
    - 在body里的raw里边写内容
        ```json
        {
            "user_name":"aoho",
            "password":"123456"
        }
        ```
- create-admin
    - api: `127.0.0.1:9009/user/admin/create`
    - method: `POST`
    - 在body里的raw里边写内容
        ```json
        {
            "user_name":"admin",
            "password":"admin",
            "age":18 
        }
        ```
- check-admin
    - api: `127.0.0.1:9009/user/admin/check`
    - method: `POST`
    - 在body里的raw里边写内容
        ```json
        {
            "user_name":"admin",
            "password":"admin"
        }
        ```

---
- Seckill
    - api: `127.0.0.1:9031/sec/kill`
    - method: `POST`
    - 在body里的raw里边写入内容
        ```json
        {
            "product_id":1,
            "user_id":114514,
            "source":"192.168.0.1",
            "auth_code":"auth_code",
            "sec_time": {{$timestamp}},
            "client_refence":"test"
        }
        ```

- SecList
    - api: `127.0.0.1:9031/sec/list`
    - method: `GET`

- SecInfo
    - api: `127.0.0.1:9031/sec/info`
    - method: `POST`
    - 在body里的raw里边写入内容
        ```json
        {
            "product_id":1
        }
        ```

---
- oauth-token：
    - api: `127.0.0.1:9019/oauth/token?grant_type=password`
    - method: `POST`
    - Authorization: 
        - Type: Basic Auth
        - Username：填写clientId
        - Password：填写clientSecret
    - 在body里的raw里边填写信息
        - username：登录账号
        - password：登录密码
        - id_type: 0是普通用户、1是管理员
- 

---
## 笔记

- gRPC流程：
    - pb.request通过decodeGRPCRequest转换为endpoint.request
    - endpoint.request通过Endpoint处理，返回endpoint.response
    - endpoint.response通过encodeGRPCResponse转换为pb.response

---
- 登录鉴权流程（oauth-token）：
    - oauth-service层
        1. 访问 `127.0.0.1:9019/oauth/token?grant_type=password` 接口
        2. 执行 **makeClientAuthorizationContext()** 函数，获取clientId和clientSecret内容
            1. 调用service层中的 **GetClientDetailByClientId()** 函数获取clientDetails，出现错误则在context中设置错误信息
            2. service层的函数调用model层的 **GetClientDetailByClientId()** 函数，从数据库中根据clientId查询ClientDetails，判断和clientSecret是否匹配，不匹配返回错误
        3. 进入**decodeTokenRequest**，从参数中获取grant_type，组装一个endpoint.TokenRequest返回
        4. 进行TokenEndpoint处理请求，它调用service.TokenGranter的 **Grant()** 方法进行授权
            1. **Grant()** 方法根据grant_type选择一个授权器（默认password），然后调用**userDetailsService.GetUserDetailByUsername()**
            2. **GetUserDetailByUsername()** 方法调用 **/pkg/client** 模块下的**userClient.CheckUser()** 方法，传入pb.UserRequest
            3. **CheckUser()** 方法里调用 **DecoratorInvoke()** 方法，它使用断路器封装了起来，内部执行服务发现，负载均衡，获取远程rpc端口，然后执行对 **user-service** 层执行rpc调用，获取pb.UserResponse

                - user-service层
                    1. rpc远程调用请求进入 **user-service** 层的 **DecodeGRPCUserRequest** 里，将pb.UserRequest转化为endpoint.UserRequest，进入 **UserEndpoint**
                    2. **UserEndpoint**调用**service**层的 **Check()** 方法处理，将返回结果封装为 **endpoint.UserResponse** 返回
                    3. service层的 **Check()** 方法调用model层的 **CheckUser()** 方法，从数据库中根据username查询用户是否存在
                    4. **endpoint.UserResponse**进入 **EncodeGRPCUserResponse()** ，被转换为**pb.UserResponse**
            4. 将pb.UserResponse组装成UserDetails，作为返回GetUserDetailByUsername()的返回值
            5. 上一步的结果返回Grant()方法中，成为UserDetails

            6. 将clientDetails和UserDetails组装成OAuth2Details，调用**tokenService.CreateAccessToken()**，获取OAuth2Token
        5. **Grant()** 授权方法返回OAuth2Token，组装成endpoint.TokenResponse返回
        6. 进入**encodeJsonResponse**，将endpoint.TokenResponse编码，返回

---
- 商品秒杀流程
    - sk-app层
        1. 访问`127.0.0.1:9031/sec/kill`接口
        2. 进入decodeSecKillRequest，从r.Body中解析出model.SecRequest
        3. 进入SecKillEndpoint处理，调用service层的SecKill(*model.SecRequest)处理
        4. SecKill先调用AntiSpam函数，判断购买者ip、id是否被封禁等，进行防作弊处理
        5. SecKill再调用SecInfoById(productId)函数，判断商品是否还在销售
        6. SecKill将请求推入**config.SkAppContext.SecReqChan**这个channel中，并启动定时器
        7. WriteHandle函数从channel中读出请求，将请求放入**conf.Redis.Proxy2layerQueueName**这个redis队列中

            - sk-core层
                1. HandleReader从**conf.Redis.Proxy2layerQueueName**这个redis队列中读取请求数据
                2. HandleReader将请求放入**config.SecLayerCtx.Read2HandleChan**这个channel中
                3. HandleUser从**config.SecLayerCtx.Read2HandleChan**读取请求，交给HandleSeckill函数处理，等待结果
                4. HandleSeckill函数进行一系列逻辑判断后，返回一个SecResult所作为结果
                5. HandleUser获取处理结果后，将结果放入**config.SecLayerCtx.Handle2WriteChan**这个channel返回
                6. HandleWrite从**config.SecLayerCtx.Handle2WriteChan**这个channel中读取结果
                7. HandleWrite将结果推入**conf.Redis.Layer2proxyQueueName**这个redis队列中
         
        8. ReadHandle函数从**conf.Redis.Layer2proxyQueueName**这个redis队列中获取请求结果，将结果放入userKey对应的channel中
        9. SecKill函数从SecResult的channel中获取请求结果，将结果返回给SecKillEndpoint
       10. 返回结果进入encodeResponse，进行编码，返回响应

---

## TODO
activity.proto内的SecProductInfoConf待修改

前端页面图片等数据，从服务器获取