local val = redis.call("get",KEYS[1])
-- 空值 这里lua 脚本会解析成 false
if not val then
    --    key 不存在
    return redis.call("set",KEYS[1],ARGV[1],'PX',ARGV[2])
elseif val == ARGV[1] then
    -- 上一次加锁成功 刷新过期时间
    redis.call('expire',KEYS[1],ARGV[2])
    return "OK"
else
    -- 被人持有锁
    return ""
end