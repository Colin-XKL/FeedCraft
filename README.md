![logo.png](asset/logo-header.png)

# Feed Craft

Craft all your feed in one place!

![author](https://img.shields.io/badge/author-Colin-blue)
![GitHub License](https://img.shields.io/github/license/Colin-XKL/FeedCraft)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/Colin-XKL/FeedCraft/docker-publish.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Colin-XKL/FeedCraft)
![GitHub Release](https://img.shields.io/github/v/release/Colin-XKL/FeedCraft)
![GitHub Discussions](https://img.shields.io/github/discussions/Colin-XKL/FeedCraft?link=https%3A%2F%2Fgithub.com%2FColin-XKL%2FFeedCraft%2Fdiscussions)
![Container Image Size On GHCR](https://ghcr-badge.egpl.dev/colin-xkl/feed-craft/size?color=%2344cc11&tag=latest&label=image+size)
![Static Badge](https://img.shields.io/badge/arch-amd64|arm64|armv7-blue)
![Open Source Love](https://badges.frapsoft.com/os/v2/open-source.svg?v=103)

Doc: English(WIP) | 简体中文 | 繁体中文(WIP)

**FeedCraft** is a powerful tool to process your rss feeds as a middleware, use it to translate your feed, extract
fulltext, emulate browser
to render js-heavy page, use llm such as google gemini to generate brief for your rss article, use natural language to
filter your rss feed, and more!

**FeedCraft** 是一个简单、强大的RSS 源处理工具,他可以作为一个中间件处理你的RSS源.
你可以用它来翻译、提取正文、模拟浏览器来渲染那些动态生成的网页并提取全文、通过大语言模型如Google
Gemini来生成文章摘要、通过自然语言筛选文章等

## 快速开始

访问以下URL 即可快速调用FeedCraft对输入的RSS源进行指定的处理
/craft/{choose_craft_option}?input_url={input_rss_url}

FeedCraft中的几个核心概念:

- Craft(工艺),指要如何处理一个rss源, 比如是要进行翻译,还是提取正文,还是AI生成摘要等

目前可用的几个选项(craft option):

- proxy: 简易RSS代理, 不作任何处理
- fulltext: 获取全文
- fulltext-plus: 获取全文,但是会模拟浏览器渲染网页,适用于常规模式无法获取到文章内容,动态渲染内容的站点
- introduction: 调用AI为文章生成摘要,附加在原文开头

你可以使用提供的demo站点快速开始体验 :
https://feed-craft.colinx.one

*注意:Demo站点仅供体验使用

## 部署

你可以通过docker快速自行部署一个FeedCraft实例,以获得更好的使用体验.
下面为一个最小docker compose 示例:

```yaml
version: "3"
services:
  app.feedcraft:
    image: ghcr.io/colin-xkl/feed-craft
    container_name: feedcraft
    restart: always
    ports:
      - 10088:8080  # 10088可替换为任何你想使用的端口
    environment:
      FC_PUPPETEER_HTTP_ENDPOINT: http://service.browserless:3000 # 替换为你自己的 browserless 或其他浏览器实例地址
      FC_REDIS_URI: redis://service.redis:6379/ # 替换为你自己的redis 实例地址
      FC_GEMINI_SECRET_KEY: <your-google-gemini-key-here> # 使用你自己的google gemini key
```

你也可以直接在一个compose文件中把redis等附加组件也一起部署好:

```yaml
version: "3"
services:
  app.feedcraft:
    image: ghcr.io/colin-xkl/feed-craft
    container_name: feedcraft
    restart: always
    ports:
      - 10088:8080  # 10088可替换为任何你想使用的端口
    environment:
      FC_PUPPETEER_HTTP_ENDPOINT: http://service.browserless:3000 # 替换为你自己的 browserless 或其他浏览器实例地址
      FC_REDIS_URI: redis://service.redis:6379/ # 替换为你自己的redis 实例地址
      FC_GEMINI_SECRET_KEY: <your-google-gemini-key-here> # 使用你自己的google gemini key
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

## 关于FeedCraft

FeedCraft 的名称和Logo参考并致敬两款游戏: MineCraft和塞尔达, 初衷和愿景是做一个简单易用、同时足够灵活, 能够有更多可能性的RSS工具.
使用问题、建议等欢迎在Discussion区讨论交流


## RoadMap

- [x] common openai api calling 
- [x] translate article and title
- [ ] feed combination support
- [ ] feed limit support
- [ ] feed custom keyword filter support
- [ ] feed natural language filter support


## 许可

GPLv3
