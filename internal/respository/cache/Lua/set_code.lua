-- 验证码在Redis上的key
local key = KEYS[1]
-- 验证次数，我们一个验证码，最多重复三次，这个用来记录重复了几次
local cntKey = key..":cnt"
-- 你的验证码
local val = ARGV[1]
-- 过期时间
local ttl = redis.call("ttl",key)
if ttl == -1 then
    -- key 存在，但是没有过期时间
    return -2
    -- -2是key不存在，ttl<540是发一个验证码已经过去一分钟了
elseif ttl == -2 or ttl < 540 then
    redis.call("set",key,val)
    redis.call("expire", key, 600)
    redis.call("set",cntKey,3)
    redis.call("expire",cntKey,600)
    --完美符合预期
    return 0
else
    --发送了一个验证码，但是还不到过期时间
    return -1
end