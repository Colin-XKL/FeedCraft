![logo.png](asset/logo-header.png)

# FeedCraft

Craft all your feed in one place!

![author](https://img.shields.io/badge/author-Colin-blue)
![GitHub License](https://img.shields.io/github/license/Colin-XKL/FeedCraft)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/Colin-XKL/FeedCraft/docker-publish.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Colin-XKL/FeedCraft)
![GitHub Release](https://img.shields.io/github/v/release/Colin-XKL/FeedCraft)
![GitHub Discussions](https://img.shields.io/github/discussions/Colin-XKL/FeedCraft?link=https%3A%2F%2Fgithub.com%2FColin-XKL%2FFeedCraft%2Fdiscussions)
![Container Image Size On GHCR](https://ghcr-badge.egpl.dev/colin-xkl/feed-craft/size?color=%2344cc11&tag=latest&label=image+size)
![Arch Badge](https://img.shields.io/badge/arch-amd64|arm64-blue)
![Open Source Love](https://badges.frapsoft.com/os/v2/open-source.svg?v=103)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/Colin-XKL/FeedCraft)

Doc: [English](https://feed-craft-doc.vercel.app/en) | [简体中文](https://feed-craft-doc.vercel.app/zh)

**FeedCraft** is a simple, powerful RSS feed processing tool.  
It can serve as middleware to handle your RSS feeds; you can use it to extract the main text, perform intelligent translation, generate article summaries through AI, filter articles using natural language, and more.  
It includes a built-in visual RSS generator (HTML/Curl/Search to RSS) that can turn web pages, API responses, or search results into RSS feeds.

**FeedCraft** 是一个简单、强大的 RSS 源处理工具.
他可以作为一个中间件处理你的RSS源, 你可以用它来提取正文、智能翻译、通过AI生成文章摘要、通过自然语言筛选文章等.
它内置了可视化 RSS 生成器 (HTML/Curl/Search to RSS)，支持将网页、API (Curl) 或搜索结果转换为 RSS 订阅源。

## 核心特性

- 开源可自部署. 可以作为中间件与现有的任何RSS阅读器协同使用
- AI Power, 可以接入Open AI接口兼容的LLM对RSS进行处理, 可自定义prompt
- **HTML/Curl/Search to RSS**: 内置可视化 RSS 生成器，支持将网页、API (Curl) 或搜索结果转换为 RSS 订阅源
- 支持**便携模式**(portable mode, 即用即走, 只需要在原RSS地址前面加个前缀即可), 和**高级模式**(dock mode,
  在后台页面自定义RSS地址和各类深度加工参数)

## 快速开始

访问以下URL格式即可快速调用FeedCraft对输入的RSS源进行指定的处理
`https://feed-craft.colinx.one/craft/{craft_name_here}?input_url={input_rss_url}`

FeedCraft中的几个核心概念:

- AtomCraft(原子工艺), 指要如何处理一个rss源, 比如是要进行翻译,还是提取正文,还是AI生成摘要等
- FlowCraft(组合工艺), 多个 AtomCraft 组成的序列, 比如你可以定义一个新的名叫 clean-article 的 FlowCraft,
  将提取全文、AI筛选文章、AI摘要组合到一起,
- Recipe(配方), 记录了以什么样的 Craft 或 FlowCraft 对某个指定的rss源进行处理, 比如你可以指定一个名叫my-zhihu-daily的recipe,
  对知乎日报的rss自动进行AI生成摘要的操作, 这个recipe对应一个新的rss地址, 你可以直接订阅这个地址得到带摘要版本的知乎日报

你可以先开始尝试下面的几个原子工艺(AtomCraft):

- **proxy**: 简易RSS代理, 不作任何处理
- **limit**: 限制文章个数, 默认最新10个
- **fulltext**: 获取全文
- **fulltext-plus**: 获取全文,但是会模拟浏览器渲染网页,适用于常规模式无法获取到文章内容,动态渲染内容的站点
- **introduction**: 调用AI为文章生成摘要,附加在原文开头
- **summary**: 调用AI总结文章主要内容,附加在原文开头
- **translate-title**: 调用AI翻译文章标题
- **translate-content**: 调用AI翻译文章内容, 只输出翻译后内容
- **translate-content-immersive**: 调用AI翻译文章内容, 沉浸式翻译模式, 每个原文段落后面添加翻译后内容
- **beautify-content**: 调用AI美化文章排版，去除广告和无关信息
- **ignore-advertorial**: 调用AI对文章进行筛选, 排除营销软文

你可以使用提供的demo站点快速开始体验 :
https://feed-craft.colinx.one

例如, 现在你想要阅读这个RSS源 `https://feeds.feedburner.com/visualcapitalist`, 但是由于是全英文的不方便快速略读,
使用FeedCraft对标题进行翻译,
那么只需要将订阅地址改成
`https://feed-craft.colinx.one/craft/translate-title?input_url=https://feeds.feedburner.com/visualcapitalist` 即可.

(部分RSS订阅软件不会自动进行字符转义, 你可以借助FeedCraft 后台的工具帮你自动进行转义)

你可以使用下面几个 RSS 快速进行测试

- 全英文，有全文
  https://feeds.feedburner.com/visualcapitalist
- 全英文，无全文
  https://ourworldindata.org/atom.xml

\*注意: Demo站点仅供体验使用

## 部署

你可以通过docker快速自行部署一个FeedCraft实例,以获得更好的使用体验.
下面为一个最小docker compose 示例:

控制台默认账号`admin`, 密码 `adminadmin`, 登陆后请尽快修改默认密码

```yaml
version: "3"
services:
  app.feed-craft:
    image: ghcr.io/colin-xkl/feed-craft
    # image: colinxkl/feed-craft  # 你也可以使用docker hub的镜像
    container_name: feed-craft
    restart: always
    ports:
      - "10088:80" # 10088可替换为任何你想使用的端口
    volumes:
      - ./feed-craft-db:/usr/local/feed-craft/db # db file
    environment:
      FC_PUPPETEER_HTTP_ENDPOINT: http://service.browserless:3000 # 替换为你自己的 browserless 或其他浏览器实例地址
      FC_REDIS_URI: redis://service.redis:6379/ # 替换为你自己的redis 实例地址
      FC_LLM_API_BASE: https://xxxxxx # LLM API 接口路径，需要以 “/v1” 结尾
      FC_LLM_API_KEY: skxxxxxx # 鉴权的key
      FC_LLM_API_MODEL: gemini-pro/chatgpt-3.5/... # 默认使用的模型
      FC_LLM_API_TYPE: openai # openai 或 ollama
      FC_DEFAULT_TARGET_LANG: zh-CN # (Optional) LLM 处理任务的默认目标语言，如 zh-CN, en-US
```

你也可以直接在一个compose文件中把redis等附加组件也一起部署好:

```yaml
version: "3"
services:
  app.feed-craft:
    image: ghcr.io/colin-xkl/feed-craft
    container_name: feed-craft
    restart: always
    ports:
      - "10088:80" # 10088可替换为任何你想使用的端口
    volumes:
      - ./feed-craft-db:/usr/local/feed-craft/db # db file
    environment:
      FC_PUPPETEER_HTTP_ENDPOINT: http://service.browserless:3000 # 替换为你自己的 browserless 或其他浏览器实例地址
      FC_REDIS_URI: redis://service.redis:6379/ # 替换为你自己的redis 实例地址
      FC_LLM_API_BASE: https://xxxxxx # LLM API 接口路径，需要以 “/v1” 结尾
      FC_LLM_API_KEY: skxxxxxx # 鉴权的key
      FC_LLM_API_MODEL: gemini-pro/chatgpt-3.5/... # 默认使用的模型
      FC_LLM_API_TYPE: openai # openai 或 ollama
      FC_DEFAULT_TARGET_LANG: zh-CN # (Optional) LLM 处理任务的默认目标语言
  service.redis:
    image: redis:6-alpine
    container_name: feedcraft_redis
    restart: always
  service.browserless:
    image: browserless/chrome
    container_name: feedcraft_browserless
    environment:
      USE_CHROME_STABLE: true
    restart: unless-stopped
```

## 关于 FeedCraft

FeedCraft 的名称和Logo参考并致敬两款游戏: MineCraft和塞尔达, 初衷和愿景是做一个简单易用、同时足够灵活, 能够有更多可能性的RSS工具.
使用问题、建议等欢迎在Discussion区讨论交流

## 许可

GPLv3

## 页面展示

**欢迎页**
<img width="1900" alt="Xnip2024-08-08_00-48-49" src="https://github.com/user-attachments/assets/d3541f4e-9ab4-4948-9fc7-5b815db225ce">

**RSS源比较工具**
你可以输入一个RSS 源的地址, 然后选择要使用的 AtomCraft, 就可以看到前后的对比.
如下图我应用的 Craft 是自定义的, 只显示文章标题带有`RSS`的文章, 可以看到筛选后的结果少了很多

<img width="1907" alt="Xnip2024-08-08_01-05-09" src="https://github.com/user-attachments/assets/7abd764a-8b19-4a72-8c94-e3ea442ff385">

**创建 AtomCraft 自定义Prompt**
<img width="1918" alt="Xnip2024-08-08_00-46-13" src="https://github.com/user-attachments/assets/ff15fe79-3792-4a96-b991-f121d2a8973e">

**自定义recipe, 指定使用哪个 AtomCraft 处理哪个RSS源**
<img width="1900" alt="Xnip2024-08-08_00-48-34" src="https://github.com/user-attachments/assets/5794de1d-28b6-45ff-8737-16f8adc6ed8a">
