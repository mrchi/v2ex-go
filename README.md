# v2ex-go

[V2EX](https://v2ex.com/) API 2.0 Beta 的非官方 Golang 库。

文档见 [V2EX › API 2\.0 Beta](https://v2ex.com/help/api)

## 安装

```bash
go get github.com/mrchi/v2ex-go
```

## API Client

### 创建 API Client

```go
var client = v2ex.NewClient("<token>")
```

### API 调用

创建新的令牌
```go
client.CreateToken(scope TokenScope, expiration TokenExpiration) (result createTokenResult, err error)
```

删除指定的提醒
```go
client.DeleteNotification(notificationId int) (err error)
```

获取指定节点
```go
client.GetNode(nodeName string) (result v2exNode, err error)
```

获取指定节点下的主题
```go
client.GetNodeTopics(nodeName string, page int) (result []v2exTopic, err error)
```

获取最新的提醒
```go
client.GetNotifications(page int) (result []v2exNotification, err error)
```

获取自己的 Profile
```go
client.GetSelfProfile() (result v2exSelfProfile, err error)
```

查看当前使用的令牌
```go
client.GetToken() (result v2exToken, err error)
```

获取指定主题
```go
client.GetTopic(topicID int) (result getTopicResult, err error)
```

获取指定主题下的回复
```go
client.GetTopicReplies(topicID int, page int) (result []v2exReply, err error)
```
