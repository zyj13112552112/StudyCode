交换机
交机类型: direct（路由）、topic（路由匹配）、fanout（广播）、headers（很少用了）
![img.png](img.png)

fanout交换机会将消息发送到所有跟交换机绑定的队列

direct交换机根据key转发

topic更加灵活，可以进行灵活的路由匹配
`*`表示匹配一个单词
`#`表示匹配0个或者多个单词