# 7Days go

通过模仿 + 改进 `Geektutu` 的简易框架轮子制作，学习go语言

## Project 1 : Gweb 

实现一个类似于Gin的web容器，完成对 `net/http` 的扩容达到实现更高效的web开发工具

## GWeb框架核心功能

- Routing：将请求映射到函数，支持动态路由
- Templates：使用内置模板引擎提供模板渲染机制
- Utilities：提供对Cookies，headers等处理机制
- Plugin：安装其他额外插件

需要实现+扩展的基础功能有：
- Runnable 对外启动器
- Handler 接口注册
- Context 上下文设计
- Router Trie树路由
- Group 分组控制
- Middleware 中间件
- HTML Template HTML模板
- Panic Recover 错误恢复