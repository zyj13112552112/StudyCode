轮询模式
在启动2个worker之后，会发现消息是按轮训过来了，消息会按照顺序一条给一台。
2个worker总会收到相同数量消息数(偏差1),这是因为消息模式设置的amqp.Persistent
![img.png](img.png)

