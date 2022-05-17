# 7Days go

通过模仿 + 改进 `Geektutu` 的简易框架轮子制作，学习go语言

## Project 1 : Gweb 

实现一个类似于Gin的web容器，完成对 `net/http` 的扩容达到实现更高效的web开发工具

需要实现+扩展的基础功能有：
- Runnable 对外启动器
- Handler 接口注册
- Context 上下文设计
- Router Trie树路由
- Group 分组控制
- Middleware 中间件
- HTML Template HTML模板
- Panic Recover 错误恢复