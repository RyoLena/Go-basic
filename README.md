# WebBook_git

Brief description of your project.

## Getting Started


### 2024年12-31日
1. 完成了对于分层结构DDD领域驱动设计的学习
2. 完成了通过该模型的设计，实现了对于用户的注册功能 
3. 下一步实现登录功能，后面实现Edit功能，并且将提供显示到Profile中

### 2025年1-1日
1. 完成了登录功能 
   - 用户通过邮箱和密码登录； 
   - 登录校验，可先先去看user中的login方法，然后看router中的userRouter中


### 2025年 2-21日
1. 完成了长短token的登录实现，
    - refreshToken需要使用user中的refresh路由检查 
      - 之前一直通过中间件检查出现了一直登陆不上的问题
   - 在`ioc/web`中链式调用`IgnorePathJWT`，出现初始化中断的问题
      - 解决办法就是之前是隐式返回，现在是显示返回
2. 在`service/sms`中增加了装饰器和延时检测切换供应商的功能`timeout_failover`文件
    - 前提拿到每个供应商的超时数据，记录次数，如果超过所设定的阈值则切换下一个服务商
    - 通过与运算以及下标每次都加一的方法，实现每次切换的时候不再是从0开始