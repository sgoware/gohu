# gohu

[![build](https://img.shields.io/badge/build-1.01-brightgreen)](https://github.com/StellarisW/StellarisW)[![go-version](https://img.shields.io/badge/go-~%3D1.8-30dff3?logo=go)](https://github.com/StellarisW/StellarisW)[![OpenTracing Badge](https://img.shields.io/badge/OpenTracing-enabled-blue.svg)](http://opentracing.io)[![PyPI](https://img.shields.io/badge/License-BSD_2--Clause-green.svg)](https://github.com/emirpasic/gods/blob/master/LICENSE)

> Go 实现的知乎网站后端

## 💡  简介

基于 [go-zero ](https://github.com/zeromicro/go-zero)微服务框架实现的一个简单知乎网站后端

## 🚀 功能

[演示视频](https://typora.stellaris.wang/gohu-show-video.mp4)

### 认证系统

- 颁发 `oauth2 token`  (包含 `access token` 和 `refresh token`)
- 对令牌进行解析获取认证信息与用户信息

### 用户系统

- 用户的登录与注册
- 用户改名
- 喜欢或赞同回答、评论
- 关注用户
- 收藏回答
- 获取自己的回答与评论

### 问答系统

- 对问题的增删查改
- 对回答的增删查改

### 评论系统

- 对回答进行评论
- 可以对评论进行回复

### 通知系统

- 关注的用户发布问题，回答问题

- 发布的问题被回答会得到通知
- 问题的回答被评论会得到通知
- 评论被回复会得到通知
- 回答或评论被点赞会得到通知

## 🌟 亮点

### 技术栈

![image-20220904045303624](./manifest/image/go-zero.png)

- [go-zero](https://go-zero.dev/)

> 一个集成了各种工程实践的包含 微服务框架
>
> 它具有高性能（高并发支撑，自动缓存控制，自适应熔断，负载均衡）、易拓展（支持中间件，方便扩展，面向故障编程，弹性设计）、低门槛（大量微服务治理和并发工具包、goctl自动生成代码很香），

​     同时go-zero也是现在最流行的go微服务框架，所以本项目采用go-zero为主框架搭建知乎后端

![](./manifest/image/mysql.svg)

- [mysql](https://www.mysql.com/)

> 一个关系型数据库管理系统，由瑞典MySQL AB 公司开发，属于 Oracle 旗下产品。MySQL 是最流行的关系型数据库管理系统关系型数据库管理系统之一，在 WEB 应用方面，MySQL是最好的 RDBMS (Relational Database Management System，关系数据库管理系统) 应用软件之一

​    遇事不决还得是mysql，后续有空（有钱）会	改进成 [TiDB](https://pingcap.com/zh/case/)，将其分布式化

![](./manifest/image/redis.svg)

- [redis](https://redis.io/)

> 一个开源的、使用C语言编写的、支持网络交互的、可基于内存也可持久化的Key-Value数据库

​    缓存存储还是选型最普遍的redis

<img src="./manifest/image/nsq.png" width="15%">

- [nsq](https://nsq.io/)

>  Go 语言编写的开源分布式消息队列中间件，具备非常好的性能、易用性和可用性

​    刚开始想用RabbitMQ，但是很难做到集群和水平扩展，再看了看其他的队列组件（Darner，Redis，Kestrel和Kafka），每种队列组件都有不同的消息传输保     	证，但似乎并没有一种能在可扩展性和操作简易性上都比较优秀的产品。

​    然后就看到了nsq，支持横向拓展，性能也很好，同时也是go语言原生开发的，因此选型nsq做发布订阅操作（通知系统里会用到）

<img src="./manifest/image/asynq.png" width="15%">

- [asynq](https://github.com/hibiken/asynq)

> go语言实现的高性能分布式任务队列和异步处理库，基于redis，类似sidekiq和celery

​	asynq是一个分布式延迟消息队列，本项目用它进行异步定时任务处理（缓存与数据库之间的同步操作）

<img src="./manifest/image/consul.svg" width="25%">

- [consul](https://www.consul.io/)

> 一套开源的分布式服务发现和配置管理系统，由HasiCorp公司用go语言开发的。提供了微服务系统中服务助理、配置中心、控制总线等功能

​	在consul和etcd之间比较，consul的服务发现很方便，也有健康检查，多数据中心等功能，同时也是go云原生项目，因此选型consul

<img src="./manifest/image/jaeger.svg" width="15%">

- [jaeger](https://www.jaegertracing.io/)

> 由Uber开源的分布式追踪系统

​	go-zero框架集成了对jaeger的支持，因此使用jaeger做追踪系统

<img src="./manifest/image/apollo.svg" width="10%">

- [apollo](https://www.apolloconfig.com/)

> 一款可靠的分布式配置管理中心，诞生于携程框架研发部，能够集中化管理应用不同环境、不同集群的配置，配置修改后能够实时推送到应用端，并且具备规范的权限、流程治理等特性，适用于微服务配置管理场景

​	使用apollo做配置管理系统，可以有效的在认证系统，用户系统，问答系统等等不同的环境下进行配置的管理

<img src="./manifest/image/traefik-logo.png" width="20%">

- [traefik](https://www.jaegertracing.io/)

> 一个为了让部署微服务更加便捷而诞生的现代HTTP反向代理、负载均衡工具

​	本项目写的微服务很多，用nginx难管理，也比较懒得写配置文件，所以用traefik，虽然性能没nginx好，但是对微服务的反向代理和负载均衡的支持很便捷，

​	同时使用traefik中的http中间件，oauth proxy也很方便

<img src="./manifest/image/docker.svg" width="15%">

- [docker](https://www.docker.com/)

> Google 公司推出的 Go 语言 进行开发实现，基于 Linux 内核的 cgroup，namespace，以及 AUFS 类的 Union FS 等技术的一个容器服务

​	容器用docker-compose部署

<img src="./manifest/image/drone.svg" width="20%">

- [drone](https://www.drone.io/)

> 一款基于go编写的容器技术持续集成工具，可以直接使用YAML配置文件即可完成自动化构建、测试、部署任务

​	结合gogs仓库，来对项目进行CI/CD持续集成，部署更方便了

### 后端

#### 认证系统

- 颁发令牌

> 对向OAuth2服务器发送请求颁发令牌的客户端进行鉴权, 鉴权成功后 颁发token
>
> 提高了令牌的安全性, 同时也可以根据客户端信息自定义 Oauth2 token

<!--app/service/oauth/model/token_granter.go-->

```go
clientId, clientSecret, userId, ok := parseBasicAuth(auth)
if !ok || clientSecret != tokenGranter.ClientDetails[clientId].ClientSecret {
   return nil, ErrInvalidAuthorizationRequest
}
```

解析认证头，判断 `clientId` 和 `clientSecret` 是否在存储库中

鉴权成功后向用户颁发 token

- 令牌存储

> token存储在redis中, 自动实现令牌的过期功能
>
> 加快用户的鉴权操作

![](./manifest/image/jwt-store.png)

<!--app/service/oauth/rpc/token/store/internal/logic/storetokenlogic.go-->

```go
l.svcCtx.Rdb.Set(l.ctx,
   model.JwtToken+"_"+strconv.FormatInt(in.UserId, 10),
   accessTokenString,
   time.Unix(in.AccessToken.RefreshToken.ExpiresAt, 0).Sub(time.Now()))
```

- 令牌刷新

> access token的过期时间为1天，refresh token的过期时间为7天，当access token过期后，服务器会自动使用一次性的refresh token 对access token进行刷新操作，获取新的access token。
>
> 这样会有效降低因 access token 泄露而带来的风险。

- 令牌传递

> 使用Cookie进行在用户与服务器之间的Oauth2 token传递，同时对cookie进行SHA256加盐，保证了cookie与token的安全性
>
> 在认证中间件中(app/common/middware/authMiddleware.go)对cookie进行哈希校验，防止用户篡改，校验通过后解析cookie自动获取token的元信息

<!--app/common/middleware/authMiddleware.go-->

```go
res, err := req.NewRequest().SetFormData(map[string]string{"oauth2_token": accessToken, "token_type": model.AccessToken}).
   Post("https://" + m.Domain + "/api/oauth/token/check")
if err != nil {
   logx.Errorf("%v", err)
   return
}
if res.StatusCode != http.StatusOK {
   logx.Errorf("%v", res)
   return
}
j := gjson.Parse(res.String())
ok = j.Get("ok").Bool()
if !ok {
   //不ok则认证失败，包括刷新令牌
   //认证的时候若认证令牌过期，则刷新令牌
   msg := j.Get("msg").String()
   // 认证令牌过期,用刷新令牌刷新
   if msg == "accessToken is expired" {
      ok = cookieWriter.Get("refresh-token", &refreshToken)
      if accessToken == "" || !ok {
         response.ResultWithData(w, http.StatusForbidden, "illegal access", map[string]interface{}{"reload": true})
         return
      }
      res, err = req.NewRequest().SetPathParam("refresh-token", refreshToken).
         Post("https://" + m.Domain + "/api/oauth/token/refresh")
      if err != nil {
         logx.Errorf("%v", err)
         return
      }
      if res.StatusCode != http.StatusOK {
         return
      }
      j = gjson.Parse(res.String())
      accessToken = j.Get("data.access_token.token_value").String()
      refreshToken = j.Get("data.access_token.refresh_token.token_value").String()
      cookieWriter.Set("x-token", accessToken)
      cookieWriter.Set("refresh-token", refreshToken)
   }
}
```

#### 用户系统

- 登录与注册

> 用户注册或首次登录得到令牌，可以在令牌有效期内实现自动登录操作，
>
> 在用户注册或登录后设置注册/登录缓存，同时特别针对登录设置了空缓存，防止大量无效请求造成缓存穿透

<!--app/service/user/rpc/crud/internal/logic/registerlogic.go-->

```go
ok, err := l.svcCtx.Rdb.SIsMember(l.ctx,
   "user_register_set",
   in.Username).Result()
```

<!--app/service/user/rpc/crud/internal/logic/loginlogic.go-->

```go
	// 在数据库中查找用户
	userSubjectModel := l.svcCtx.UserModel.UserSubject
	userSubject, err := userSubjectModel.WithContext(l.ctx).Where(userSubjectModel.Username.Eq(in.Username)).First()
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		// 设置空缓存,防止大量非法请求造成缓存穿透
		err = l.svcCtx.Rdb.Set(l.ctx,
			fmt.Sprintf("user_login_%s", in.Username),
			fmt.Sprintf("%d:%d", 0, 0),
			time.Second*86400).Err()
		if err != nil {
			logger.Errorf("update [user_login] cache failed, err: %v", err)
		}
		res = &pb.LoginRes{
			Code: http.StatusNotFound,
			Msg:  "uid not exist",
			Ok:   false,
		}
		logger.Debugf("send message: %v", res.String())
		return res, nil
	default:
		{
			logger.Errorf("query [user_subject] in mysql failed, err: %v", err)
			res = &pb.LoginRes{
				Code: http.StatusInternalServerError,
				Msg:  "internal err",
				Ok:   false,
			}
			logger.Debugf("send message: %v", res.String())
			return res, err
		}
	}
	
	...
	
		err = l.svcCtx.Rdb.Set(l.ctx,
		fmt.Sprintf("user_login_%s", in.Username),
		fmt.Sprintf("%d:%s", userSubject.ID, userSubject.Password),
		time.Second*86400).Err()
	if err != nil {
		logger.Errorf("set [user_login] cache failed, err: %v", err)
	}
```

- IP 归属地

> 对用户登录的IP进行解析（具体看我 [ip-parse](https://github.com/StellarisW/ip-parse) 仓库），在回答内容，评论内容中设置用户IP归属地

<!--app/utils/net/ip/ip.go-->

```go
func GetIpLocFromApi(ip string) (loc string) {
   var err error
   if domain == "" {
      domain, err = apollo.GetMainDomain()
      if err != nil {
         return "未知"
      }
   }
   apiAddr := "http://ip." + domain + "/api/parse?ip=" + ip
   res, err := req.NewRequest().Get(apiAddr)
   j := gjson.Parse(res.String())
   if j.Get("ok").Bool() == false {
      return "未知"
   }
   locStr := j.Get("location").String()
   output := strings.Split(locStr, "|")
   if output[0] != "中国" {
      return output[0]
   }
   strings.Trim(output[2], "省")
   return output[2]
}
```

- 用户信息缓存

> 对`user-subject`,`user-collect`表进行缓存，针对高频更新的字段(如 `user-subject` 中的 `follower` ， `user-collect` 中的赞同操作)进行定时更新数据库的操作

<!--app/service/mq/asynq/processor/internal/logic/user/task.go-->

```go
func (l *ScheduleUpdateUserSubjectRecordHandler) ProcessTask(ctx context.Context, _ *asynq.Task) (err error) {
    // 获取缓存中需要更新用户 follower 字段的用户id
   members, err := l.Rdb.SMembers(ctx,
      "user_follower_cnt_set").Result()
   if err != nil {
      return fmt.Errorf("get [user_follower] member failed, err: %v", err)
   }
   l.Rdb.Del(ctx,
      fmt.Sprintf("user_follower_cnt_set"))

   userSubjectModel := l.UserModel.UserSubject
   for _, member := range members {
       // 获取缓存中 follower 的数量
      followerCount, err := l.Rdb.Get(ctx,
         fmt.Sprintf("user_follower_cnt_%s", member)).Int()
      if err != nil {
         return fmt.Errorf("get [user_follower] cnt failed, err: %v", err)
      }

      err = l.Rdb.Del(ctx,
         fmt.Sprintf("user_follower_cnt_%s", member)).Err()
      if err != nil {
         return fmt.Errorf("del [user_follower] cnt failed, err: %v", err)
      }
	
       // 更新数据库
      userSubject, err := userSubjectModel.WithContext(ctx).
         Select(userSubjectModel.ID, userSubjectModel.Follower).
         Where(userSubjectModel.ID.Eq(cast.ToInt64(member))).
         First()
      if err != nil {
         return fmt.Errorf("get [user_subject] record failed, err: %v", err)
      }

      _, err = userSubjectModel.WithContext(ctx).
         Select(userSubjectModel.ID, userSubjectModel.Follower).
         Where(userSubjectModel.ID.Eq(cast.ToInt64(member))).
         Update(userSubjectModel.Follower, int(userSubject.Follower)+followerCount)
      if err != nil {
         return fmt.Errorf("update [user_subject] record failed, err: %v", err)
      }
   }

   return nil
}
```

#### 问答系统

- 问题与回答信息缓存

> 也是和上面一样的逻辑，对表进行缓存，高频字段定时数据库同步

- 回答点赞

<!--app/service/user/rpc/crud/internal/logic/docollectionlogin.go-->

```go
// 通知用户
err = notificationMqProducer.PublishNotification(producer, notificationMqProducer.PublishNotificationMessage{
   MessageType: 2,
   Data: notificationMqProducer.ApproveAndLikeData{
      UserId:  in.UserId,
      Action:  1,
      ObjType: in.ObjType,
      ObjId:   in.ObjId,
   },
})
if err != nil {
   return fmt.Errorf("publish notificaion to nsq failed, %v", err)
}

// 更新缓存
err = svcCtx.Rdb.Incr(ctx,
   fmt.Sprintf("answer_index_approve_cnt_%d", in.ObjId)).Err()
if err != nil {
   return fmt.Errorf("incr [answer_index_approve_cnt] failed, %v", err)
}

err = svcCtx.Rdb.SAdd(ctx,
   "answer_index_approve_cnt_set",
   in.ObjId).Err()
if err != nil {
   return fmt.Errorf("update [answer_index_approve_cnt_set] failed, err: %v", err)
}
```

#### 评论系统

- 低耦合

> 将评论这个模块单独从回答中抽离出来，方便后续的功能需要评论模块的使用

- 信息全面

> 显示回复的用户id和被回复的用户id，回复的ip归属地等等

#### 通知系统

- 发布订阅

> 使用nsq对用户的各种操作(如关注，点赞等等)对操作的对象(被关注的人，被点赞回答的作者)进行通知

<!--app/service/mq/nsq/consumer/internal/listen/notification/handler.go-->

```go
func (m *PublishNotificationHandler) HandleMessage(nsqMsg *nsq.Message) (err error) {
   msg := &notificationMqProducer.PublishNotificationMessage{}
   err = json.Unmarshal(nsqMsg.Body, &msg)
   if err != nil {
      return fmt.Errorf("unmarshal msg failed, %v", err)
   }

   ctx := context.Background()
   switch msg.MessageType {
   case 1:
      // 关注我的
      data := &notificationMqProducer.FollowerData{}

      bytesData, err := json.Marshal(msg.Data)
      if err != nil {
         return fmt.Errorf("marshal msg data failed, %v", err)
      }

      err = json.Unmarshal(bytesData, &data)
      if err != nil {
         return fmt.Errorf("unmarshal msg data failed, %v", err)
      }

      userInfoRes, err := req.NewRequest().Get(
         fmt.Sprintf("https://%s/api/user/profile/%s", m.Domain, cast.ToString(data.FollowerId)))
      if err != nil {
         return fmt.Errorf("query user info failed, %v", err)
      }

      j := gjson.Parse(userInfoRes.String())
      if !j.Get("ok").Bool() {
         return fmt.Errorf("query user info failed, %v", j.Get("msg").String())
      }

      rpcRes, _ := m.NotificationCrudRpcClient.PublishNotification(ctx, &crud.PublishNotificationReq{
         UserId:      data.UserId,
         MessageType: 1,
         Title:       fmt.Sprintf("用户 %s 关注了你", j.Get("data.nickname").String()),
         Content:     "",                                                              // 空
         Url:         fmt.Sprintf("https://%s/profile/%d", m.Domain, data.FollowerId), // 用户主页
      })
      if !rpcRes.Ok {
         return fmt.Errorf("publish notification failed, %v", rpcRes.Msg)
      }
      
      ...
      
}
```

### 运维

#### 服务发现

> 通过go-zero框架的原生支持，来进行服务的发现与调用操作

![image-20220902154029133](./manifest/image/consul.png)

#### 链路追踪

> 通过go-zero框架的原生支持，进行对服务的链路追踪，方便debug

 ![image-20220902154140585](./manifest/image/jaeger.png)

#### 配置管理

> 在apollo平台上进行配置的更改操作，然后go后端读取apollo的配置文件后，缓存在本地，同时进行配置的热更新，
>
> 然后使用viper来进行配置文件的读取操作

![image-20220902153815361](./manifest/image/apollo.png)

#### 网关代理

> 使用traefik对微服务进行反向代理，同时自动生成CA证书，对全部的微服务进行tls认证

![image-20220902154403897](./manifest/image/traefik.png)

#### 项目部署

> 代码push到gogs仓库

![image-20220902154503351](./manifest/image/gogs.png)

> drone通过web钩子拷贝代码

![image-20220902154545273](./manifest/image/webhook.png)

> drone进行持续集成操作

![image-20220902154620156](./manifest/image/drone.png)

> docker-compose 并发构建镜像

`build.sh`

``` sh
#!/bin/bash

export PROJECT_NAME=$1

export THREAD=$2

docker_names=('oauth-api' 'oauth-rpc-token-enhancer' 'oauth-rpc-token-store' \
'user-api' 'user-rpc-crud' 'user-rpc-info' 'user-rpc-vip' 'notification-api' \
'notification-rpc-crud' 'notification-rpc-info' 'mq-asynq-scheduler' 'mq-asynq-processor' \
'mq-nsq-consumer')

function docker_build() {
  if [ "$1" -ef "" ]; then
    return 0
  fi

  array=$(echo "$1" | tr '-' '\n')
  path='./app/service'
  for var in $array
  do
    path="${path}""/""${var}"
  done
  docker build -t "$PROJECT_NAME""_""$1" "${path}"
  return 1
}

[ -e /tmp/fd1 ] || mkfifo /tmp/fd1
exec 3<>/tmp/fd1
rm -rf /tmp/fd1

for ((i=1;i<=THREAD;i++))
do
  echo >&3
done

cd /www/site/"$PROJECT_NAME" || exit

remain_build=${#docker_names[@]}

echo "start building images, remain: ""${remain_build}"

for docker_name in ${docker_names[*]}
do
  read -r -u3
{
  docker_build "${docker_name}"
  remain_build=$(expr "${remain_build}" - 1)
  echo "build ""${docker_name}"" complete, remain: ""${remain_build}"
  echo >&3
} &
done

wait

exec 3<&-
exec 3>&-
```

## 🗼架构设计

![clouddisk架构图](./manifest/image/architecture.jpg)

## 📂 存储设计

### 表设计

#### 用户系统

##### `user_subject`

##### 记录

![](./manifest/image/user_subject_record.png)

索引

![](./manifest/image/user_subject_record_index.png)

##### `user_collection`

记录

![](./manifest/image/user_collection_record.png)

索引

![](./manifest/image/user_collection_index.png)

外键

![](./manifest/image/user_collection_foreign_key.png)

#### 问答系统

##### `question_subject`

记录

![](./manifest/image/question_subject_record.png)

索引

![](./manifest/image/question_subject_index.png)



##### `question_content`

记录

![](./manifest/image/question_content_record.png)

外键

![](./manifest/image/question_content_foreign_key.png)



##### `answer_index`

记录

![](./manifest/image/answer_index_record.png)

索引

![](./manifest/image/answer_index_index.png)

外键

![](./manifest/image/answer_index_foreign_key.png)



##### `answer_content`

记录

![](./manifest/image/answer_content_record.png)

外键![](./manifest/image/answer_content_foreign_key.png)

#### 评论系统

##### `comment_subject`

记录

![](./manifest/image/comment_subject_record.png)

索引

![](./manifest/image/comment_subject_index.png)



##### `comment_index`

记录

![](./manifest/image/comment_index_record.png)

索引

![](./manifest/image/comment_index_index.png)

外键

![](./manifest/image/comment_index_foreign_key.png)



##### `comment_content`

记录

![](./manifest/image/comment_content_record.png)

外键

![](./manifest/image/comment_content_foreign_key.png)

#### 通知系统

##### `notification_subject`

记录

![](./manifest/image/notification_subject_record.png)

索引

![](./manifest/image/notification_subuject_index.png)

外键

![](./manifest/image/notification_subject_foreign_key.png)





##### `notification_content`

记录

![](./manifest/image/notification_content_record.png)



外键

![](./manifest/image/notification_content_foreign_key.png)



### 缓存设计

#### jwt 缓存

缓存 Oauth2 Token，利用 redis 特性自动实现令牌过期功能

![](./manifest/image/jwt_cache.png)

#### 表缓存

> 表记录的数据比较多，因此采用protobug进行序列化，而不是用json，
>
> 因为在protobuf的编解码性能远远高出JSON的性能。
>
> 同时占用空间也比json小很多，大大减少了缓存表的开销

下面是 `user_subject` 表的缓存样例

<!--app/service/mq/asynq/processor/internal/logic/user/task.go-->

```go
func (l *MsgCreateUserSubjectHandler) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
   var payload job.MsgCreateUserSubjectPayload
   if err = json.Unmarshal(task.Payload(), &payload); err != nil {
      return fmt.Errorf("unmarshal [MsgCreateUserSubjectPayload] failed, err: %v", err)
   }

   userSubjectId := l.IdGenerator.NewLong()

   userSubjectModel := l.UserModel.UserSubject

   now := time.Now()

   err = userSubjectModel.WithContext(ctx).
      Create(&model.UserSubject{
         ID:         userSubjectId,
         Username:   payload.Username,
         Password:   payload.Password,
         Nickname:   payload.Nickname,
         CreateTime: now,
         UpdateTime: now,
      })
   if err != nil {
      return fmt.Errorf("create [user_subject] record failed, err: %v", err)
   }

   userSubjectProto := &pb.UserSubject{
      Id:         userSubjectId,
      Username:   payload.Username,
      Password:   payload.Password,
      Nickname:   payload.Nickname,
      CreateTime: now.String(),
      UpdateTime: now.String(),
   }

   userSubjectBytes, err := proto.Marshal(userSubjectProto)
   if err != nil {
      return fmt.Errorf("marshal [userSubjectProto] into proto failed, err: %v", err)
   }

   err = l.Rdb.Set(ctx,
      fmt.Sprintf("user_subject_%d", userSubjectId),
      userSubjectBytes,
      time.Second*86400).Err()
   if err != nil {
      return fmt.Errorf("update [user_subject] cache failed, err: %v", err)
   }

   err = l.Rdb.Set(ctx,
      fmt.Sprintf("user_login_%d", userSubjectId),
      fmt.Sprintf("%d:%s", userSubjectId, payload.Password),
      time.Second*86400).Err()
   if err != nil {
      return fmt.Errorf("update [user_login] cache failed, err: %v", err)
   }

   return nil
}
```

key 格式： `[tableName]_[primaryId]`

例如 `user_subject` 的缓存

![](./manifest/image/user_subject_cache.png)

#### 注册缓存

使用集合存储已经注册的用户名（后续集合过大时，考虑随机删除集合中的用户名减小缓存压力）

key 名称：`user_register_set`

![](./manifest/image/user_register_set.png)

#### 登录缓存

使用 key 存储用户 id 与加盐后的密码

key 格式：`user_login_[username]`

![](./manifest/image/user_login_cache.png)

#### 收藏缓存

key 格式：`user_collect_set_[userId]_[collectType]_[objType]`

![](./manifest/image/user_collect_set.png)

#### 关注用户缓存

key 格式：`user_follwer_cnt_set` `user_follower_cnt_[userId]`

缓存关注用户的数量，然后根据这个数量数据库定时统一更新，大大减小数据库压力

key 格式：`user_follower_member_set_[userId]`

缓存用户下的关注者集合，方便发布通知

![](./manifest/image/user_follower_member_set.png)

#### 通知缓存

key 格式：`notification_[userId]_[notificationType]`

缓存用户通知下的一个 id 集合

#### 发布缓存

key 格式：`question_id_user_set_[userId]` `answer_id_user_set_[userId]` `comment_id_user_set_[userId]`

缓存用户发布对象的的 id 集合

## 📖 API文档

[接口文档](https://documenter.getpostman.com/view/22490304/VUxStm2d)

## ⚙ 项目结构

<details>
<summary>展开查看</summary>
<pre>
<code>
    ├── app ----------------------------- (项目文件)
        ├── common ---------------------- (全局通用目录)
        	├── config ------------------ (获取配置文件相关)
        	├── log --------------------- (日志配置)
        	├── middleware -------------- (中间件)
        	├── model ------------------- (全局模型)
        	├── mq ---------------------- (消息队列设置)
        ├── service --------------------- (微服务)
            ├── comment ----------------- (评论系统)
            ├── mq ---------------------- (消息队列服务)
            ├── notification ------------ (通知系统)
            ├── oauth ------------------- (认证系统)
            ├── question ---------------- (问答系统)
            ├── user -------------------- (用户系统)
    ├── manifest ------------------------ (交付清单)
    	├── deploy ---------------------- (部署配置文件)
    		├── docker ------------------ (docker配置文件)
    		├── kustomize --------------- (k8s配置文件)
        ├── sql ------------------------- (mysql初始化配置文件)
    ├── utils --------------------------- (工具包)              
    ├── .drone.yml ---------------------- (drone自动构建配置文件)
    ├── docker-compose.yaml ------------- (docker-compose配置文件)
</code>
</pre>
</details>


## 🛠 环境要求

- golang 版本 ~= 1.19
- mysql 版本 ~=5.7
- redis 版本 ~=5.0

## 📌 TODO

- [x] OAuth授权服务器
- [x] 用户登录/注册
- [x] 发布问题
- [x] 回答问题
- [x] 评论回答
- [x] 获取自己的所有问题、回答、评论
- [x] 删除或修改自己发布的问题、回答、评论
- [x] 赞同
- [x] 关注
- [x] 收藏
- [x] 通知
- [ ] 热榜
- [ ] 盐选会员
- [ ] 文章
