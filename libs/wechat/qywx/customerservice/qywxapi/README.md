# 企业微信客服事件响应消息使用指南

## 概述

当特定的回调事件包含 `code` 字段，或通过 `TransServiceState` 接口变更到特定的会话状态时，会返回 `msg_code`。开发者可以以这个 `msg_code` 为凭证，调用 `SendMsgOnEvent` 接口发送相应事件场景下的消息，如客服欢迎语、客服提示语和会话结束语等。

## 支持的事件场景

| 事件场景 | 允许下发条数 | code有效期 | 支持的消息类型 | 获取code途径 |
|---------|-------------|-----------|--------------|-------------|
| 用户进入会话，用于发送客服欢迎语 | 1条 | 20秒 | 文本、菜单 | 事件回调 |
| 进入接待池，用于发送排队提示语等 | 1条 | 48小时 | 文本 | 转接会话接口 |
| 从接待池接入会话，用于发送非工作时间的提示语或超时未回复的提示语等 | 1条 | 48小时 | 文本 | 事件回调、转接会话接口 |
| 结束会话，用于发送结束会话提示语或满意度评价等 | 1条 | 20秒 | 文本、菜单 | 事件回调、转接会话接口 |

## 使用步骤

### 1. 转换会话状态获取 msg_code

```go
// 将会话转入待接入池
req := &qywxapi.TransServiceStateReq{
    OpenKfID:       "wkxxxxxxxxxxxx",
    ExternalUserID: "wmxxxxxxxxxxxx",
    ServiceState:   2, // 转入待接入池排队中
}

msgCode, err := api.TransServiceState(accessToken, req)
if err != nil {
    // 处理错误
}
```

### 2. 发送文本事件响应消息

```go
// 创建文本消息
textMsg := qywxapi.NewEventTextMessage(
    msgCode,
    "欢迎咨询，请问有什么可以帮助您的？",
)

// 发送消息
msgID, err := api.SendMsgOnEvent(accessToken, textMsg)
if err != nil {
    // 处理错误
}
```

### 3. 发送菜单事件响应消息

```go
// 创建菜单消息（满意度评价）
menuMsg := qywxapi.NewEventMenuMessage(
    msgCode,
    &qywxapi.MenuContent{
        HeadContent: "您对本次服务是否满意呢？",
        List: []qywxapi.MenuItem{
            {
                Type: "click",
                Click: &qywxapi.MenuClickItem{
                    ID:      "101",
                    Content: "满意",
                },
            },
            {
                Type: "click",
                Click: &qywxapi.MenuClickItem{
                    ID:      "102",
                    Content: "不满意",
                },
            },
            {
                Type: "view",
                View: &qywxapi.MenuViewItem{
                    URL:     "https://example.com/feedback",
                    Content: "点击跳转到反馈页面",
                },
            },
            {
                Type: "miniprogram",
                MiniProgram: &qywxapi.MenuMiniProgramItem{
                    AppID:    "wx1234567890abcdef",
                    PagePath: "pages/feedback/index",
                    Content:  "打开小程序反馈",
                },
            },
        },
        TailContent: "感谢您的评价",
    },
)

// 发送消息
msgID, err := api.SendMsgOnEvent(accessToken, menuMsg)
if err != nil {
    // 处理错误
}
```

## 完整示例流程

### 场景1: 用户进入会话发送欢迎语

```go
// 1. 从事件回调中获取 welcome_code
// (在事件处理函数中)
// welcomeCode := event.WelcomeCode

// 2. 发送欢迎语
textMsg := qywxapi.NewEventTextMessage(welcomeCode, "您好，欢迎咨询！")
msgID, err := api.SendMsgOnEvent(accessToken, textMsg)
```

### 场景2: 转入待接入池发送排队提示

```go
// 1. 将会话转入待接入池
req := &qywxapi.TransServiceStateReq{
    OpenKfID:       "wkxxxxxxxxxxxx",
    ExternalUserID: "wmxxxxxxxxxxxx",
    ServiceState:   2, // 待接入池排队中
}

msgCode, err := api.TransServiceState(accessToken, req)
if err != nil {
    return err
}

// 2. 发送排队提示语
textMsg := qywxapi.NewEventTextMessage(msgCode, "您好，正在为您分配客服，请稍候...")
msgID, err := api.SendMsgOnEvent(accessToken, textMsg)
```

### 场景3: 结束会话发送满意度评价

```go
// 1. 结束会话
req := &qywxapi.TransServiceStateReq{
    OpenKfID:       "wkxxxxxxxxxxxx",
    ExternalUserID: "wmxxxxxxxxxxxx",
    ServiceState:   4, // 已结束
}

msgCode, err := api.TransServiceState(accessToken, req)
if err != nil {
    return err
}

// 2. 发送满意度评价菜单
menuMsg := qywxapi.NewEventMenuMessage(
    msgCode,
    &qywxapi.MenuContent{
        HeadContent: "您对本次服务是否满意呢？",
        List: []qywxapi.MenuItem{
            {
                Type: "click",
                Click: &qywxapi.MenuClickItem{
                    ID:      "101",
                    Content: "满意",
                },
            },
            {
                Type: "click",
                Click: &qywxapi.MenuClickItem{
                    ID:      "102",
                    Content: "不满意",
                },
            },
        },
        TailContent: "感谢您的反馈",
    },
)

msgID, err := api.SendMsgOnEvent(accessToken, menuMsg)
```

## 注意事项

1. **code只能使用一次**：每个 `msg_code` 只能使用一次，使用后立即失效。
2. **有效期限制**：不同场景的 `code` 有效期不同，请参考上表。
3. **消息类型限制**：不同场景支持的消息类型不同，发送不支持的消息类型会报错。
4. **会话状态限制**：除"用户进入会话事件"以外，响应消息仅支持会话处于获取该 `code` 的会话状态时发送。
5. **48小时内未收过欢迎语**：用户在过去48小时里未收过欢迎语，且未向客服发过消息，才能发送欢迎语。

## 参考文档

- [发送欢迎语等事件响应消息](https://developer.work.weixin.qq.com/document/path/95122)
- [分配客服会话](https://developer.work.weixin.qq.com/document/path/94669)