```
1.一致性hash，保证同一个信箱的消息，在同一个msgstore模块入库
2.群消息入队，保证同一个群的消息，被同一个goroutine入队，否则会存在消息丢失问题。【a消息先来，但是线程卡了，这时候b消息成功入库并被通知拉取了，然后a消息成功入库了，那么a消息将永远无法被拉到】
3.其实单聊也会出现群聊的那种情况，但是出现的概率及低，可以暂时不用走队列
4.写信箱的时候增加cursor缓存，避免数据库空拉取【也需要端上拉消息的时候，end用notify的msgid，而不是正无穷，某些场景下可以正无穷，但是不能频繁】
5.消息入库，增加限速保护【建议直接在接口层增加】
6.对于群消息的读取进度，每个人对每个群增加一个msg_progress表，记录读取的群消息id
```