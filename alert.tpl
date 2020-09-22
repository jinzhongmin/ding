
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