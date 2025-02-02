local key = KEYS[1]
-- 用户输入的code
local expectedCode = ARGV[1]
local cntKey = key..":cnt"

local cnt = redis.call("get",cntKey)
local code = redis.call("get",key)

if not code then
    redis.call("del", cntKey)
    return -1  --验证码不存在
end

-- 验证次数已经耗尽
if not cnt or tonumber(cnt) <= 0 then
    -- 用户已经用完了验证次数或者 cntKey 不存在
    redis.call("del", key)
    redis.call("del", cntKey)
    return -2
end
--验证码相等
--不能删除验证码，因为如果你删除了就有可能存在其他的问题
--立刻再次再次发送一个验证码
if code == expectedCode then
    -- 把次数标记为 -1，认为不可用
    redis.call("set", cntKey, -1)
    return 0
else
    -- 用户输入错误
    -- 可验证次数减一
    redis.call("decr", cntKey)
    return -1 -- 验证失败，但还有机会
end

