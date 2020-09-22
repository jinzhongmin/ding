# ding

prometheus 钉钉告警机器人

## 添加钉钉机器人

参考 https://www.dingtalk.com/qidian/help-detail-20781541.htmlhttps://www.dingtalk.com/qidian/help-detail-20781541.html

## 配置

### 程序配置

./ding.yml

``` yam
global:
  # routePath 请与alertmanager中url的路径配置一致
  routePath: '/'
  # listenPort 本地监听的端口，请与alertmanager中url的端口配置一致
  listenPort: '5001'

ding:
  # token 钉钉机器人的token
  token: 'xxxxxxxxx'
  # title 钉钉告警通知的标题, 正文不显示
  title: '告警通知'
  msgFormat: 
    # msgType text, markdown
    msgType: 'markdown'
    # msgHead 每条通知的头部显示, 可自定义 
    msgHead: '来自 Prometheus 的告警通知\n++++++\n\n'
    # msgBodyTpl 告警的显示模板, 可自定义
    msgBodyTpl: './alert.tpl'
```

### 告警模板配置

./alert.tpl

变量语法参考 golang text/template , 主要的变量默认模板里已经写了

``` tpl

{{ if eq .Alert.Status "firing"}}**告警**{{ else if eq .Alert.Status "resolved" }}**清除**{{ end }}
> 
> 告警级别 : {{ .Alert.Labels.severity }}
> 
> 告警名称 : {{ .Alert.Labels.alertname }}
>
> 告警源 : {{ .Alert.Labels.instance }}
> 
> 告警发生时间 : {{ .Alert.StartsAt.Local.Format "2006-01-02 15:04:05" }}
> 
> 告警结束时间 : {{ .Alert.EndsAt.Format "2006-01-02 15:04:05" }}
> 
> 告警描述 : {{ .Alert.Annotations.description }}
>
> 告警详情 : {{ .Alert.Annotations.summary }}
>
---
```
