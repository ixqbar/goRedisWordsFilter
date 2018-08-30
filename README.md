
## 脏字过滤

### 版本
``` 
v0.0.1
```

### 使用
``` 
./bin/goRedisWordsFilter_linux --config=config.xml

redis-cli -p 8299 --raw
127.0.0.1:8299> total
2
127.0.0.1:8299> exists 好屌
0
127.0.0.1:8299> add 好屌
OK
127.0.0.1:8299> total
3
127.0.0.1:8299> filter 妈的-我看他说话的语气，好屌啊
**-我看他说话的语气，**啊
127.0.0.1:8299> delete 妈的
OK
127.0.0.1:8299> total
2
127.0.0.1:8299> filter 妈的-我看他说话的语气，好屌啊
妈的-我看他说话的语气，**啊
127.0.0.1:8299> reload
OK
127.0.0.1:8299> total
3
127.0.0.1:8299> exists 好屌
1
127.0.0.1:8299> filter 妈的-我看他说话的语气，好屌啊
**-我看他说话的语气，**啊
127.0.0.1:8299> FLUSHALL
OK
127.0.0.1:8299> total
0
127.0.0.1:8299> reload
OK
127.0.0.1:8299> total
3
127.0.0.1:8299> filter 妈的-我看他说话的语气，好屌啊
**-我看他说话的语气，**啊
```

### 配置 config.xml
``` 
<?xml version="1.0" encoding="UTF-8" ?>
<config>
    <server>0.0.0.0:8299</server>
    <dict>dict/words.conf</dict>
</config>
```

### 交流
 * qq群 233415606