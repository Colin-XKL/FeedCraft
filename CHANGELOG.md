# Change Log

## [3.2.0](https://github.com/Colin-XKL/FeedCraft/compare/v3.1.0...v3.2.0) (2026-08-30)


### Features

* add  better error handling & observability ([68c1987](https://github.com/Colin-XKL/FeedCraft/commit/68c19871d5fd11a3092ac472a3ac1563f838b7c2))
* add admin page for inbox and system token management ([29699b5](https://github.com/Colin-XKL/FeedCraft/commit/29699b5057edf8a61833c867775df6e1903d2ff3))
* add AI content process craft ([#766](https://github.com/Colin-XKL/FeedCraft/issues/766)) ([648049b](https://github.com/Colin-XKL/FeedCraft/commit/648049b02ad041d447e0d42ae0a683db612c4312))
* add ai filter atom craft ([#753](https://github.com/Colin-XKL/FeedCraft/issues/753)) ([2b5a5c7](https://github.com/Colin-XKL/FeedCraft/commit/2b5a5c77ed9ab7b3c32257fbdd883177c2c0f5d0))
* add builder for topicFeed ([aff43de](https://github.com/Colin-XKL/FeedCraft/commit/aff43de7ae2f7f92e52246d55d6ec36892383f30))
* add cache status to craft log output ([cc52fe1](https://github.com/Colin-XKL/FeedCraft/commit/cc52fe11a6ea7978ea86e212fd0e599a320d8bed))
* add conecpt visualizer page ([a89dc27](https://github.com/Colin-XKL/FeedCraft/commit/a89dc276df84b023cfc71800f93bb9281fb1fd62))
* add direct inbox RSS endpoint GET /inbox/:inbox_id/rss ([698e779](https://github.com/Colin-XKL/FeedCraft/commit/698e779ef8f9da99a8b4f9eb2df31501cfa24e67))
* add example RSS feeds ([#755](https://github.com/Colin-XKL/FeedCraft/issues/755)) ([27cc2ac](https://github.com/Colin-XKL/FeedCraft/commit/27cc2ac86ab9988ed90557aef894f4efbe9e97ad))
* add ID format validation to controllers ([1b6a52f](https://github.com/Colin-XKL/FeedCraft/commit/1b6a52fa2c80ab176fcb0f44de32be42399c8c04))
* add inbox support ([e6fee58](https://github.com/Colin-XKL/FeedCraft/commit/e6fee58bcac9dd0bbed9ed1691fc5522291e0917))
* add io.LimitReader to prevent OOM in HTTP fetcher ([#717](https://github.com/Colin-XKL/FeedCraft/issues/717)) ([d9f2ade](https://github.com/Colin-XKL/FeedCraft/commit/d9f2ade1cad5e010530629accc0a3430380fab67))
* add link flatten atom craft ([#825](https://github.com/Colin-XKL/FeedCraft/issues/825)) ([a317405](https://github.com/Colin-XKL/FeedCraft/commit/a317405082a3a2a5d192f1c6e2b22d3979b3851c))
* add more RSS format examples ([#758](https://github.com/Colin-XKL/FeedCraft/issues/758)) ([324e525](https://github.com/Colin-XKL/FeedCraft/commit/324e525eea7431ae3f66f88ecffb3ed7b5e84abe))
* add template options to curl-to-rss ([11e6519](https://github.com/Colin-XKL/FeedCraft/commit/11e6519d169cf0d449cf07e409b44cc2407613ee))
* add topic feed ([0888588](https://github.com/Colin-XKL/FeedCraft/commit/088858887c81c4fb3ac0a6eba7928fadff6b2116))
* add topic feed aggregation and processors ([07b662a](https://github.com/Colin-XKL/FeedCraft/commit/07b662a2356045da25248c4ef03e96eebbb1e1d1))
* add web admin page support for web monitor feature ([c8717f7](https://github.com/Colin-XKL/FeedCraft/commit/c8717f791c28ca33162cee2d532fde40ad681c6d))
* **admin:** highlight missing crafts on health dashboard ([#709](https://github.com/Colin-XKL/FeedCraft/issues/709)) ([24bf35a](https://github.com/Colin-XKL/FeedCraft/commit/24bf35af3938dcf1e10a8e9964cbdb7468099c3a))
* **admin:** sort atom craft templates alphabetically ([3cac037](https://github.com/Colin-XKL/FeedCraft/commit/3cac03766b7988b888d24ce28562ffc98d508f94))
* better error message and ux for feed compare page ([0d2eb76](https://github.com/Colin-XKL/FeedCraft/commit/0d2eb7636c83c4a4fc14834492e91a84b86e2296))
* better ux for feed preview tool ([a8033a1](https://github.com/Colin-XKL/FeedCraft/commit/a8033a132e4abf0367e118fc33f67b6e084fa3a1))
* complete inbox i18n coverage and add documentation (en/zh/zh-TW) ([#779](https://github.com/Colin-XKL/FeedCraft/issues/779)) ([f9e60e9](https://github.com/Colin-XKL/FeedCraft/commit/f9e60e968179c24cf67d1fe623036850546e32f3))
* **craft:** add llm article processors and payload builder ([003541b](https://github.com/Colin-XKL/FeedCraft/commit/003541b4a6af7dc15805010ee89368e96ab8d541))
* **craft:** 调用 LLM 前的 cleanup 移除 base64 编码图片以节省 token ([#827](https://github.com/Colin-XKL/FeedCraft/issues/827)) ([9424cca](https://github.com/Colin-XKL/FeedCraft/commit/9424cca243f59cc31d5f8a9bd4be48ed42b96c18))
* **custom_recipe:** add copy link support ([df86ae2](https://github.com/Colin-XKL/FeedCraft/commit/df86ae2959f79026b66a953add102a198570775f))
* **doc-site:** 文档目录重构为四象限结构 (COL-38) ([#913](https://github.com/Colin-XKL/FeedCraft/issues/913)) ([44b7da1](https://github.com/Colin-XKL/FeedCraft/commit/44b7da10297b9049b78b727eef453dd38d6561b9))
* **doc-site:** 文档目录重构为四象限结构 (COL-38) ([#913](https://github.com/Colin-XKL/FeedCraft/issues/913)) ([#916](https://github.com/Colin-XKL/FeedCraft/issues/916)) ([e1a7d6d](https://github.com/Colin-XKL/FeedCraft/commit/e1a7d6d089e97e1b712ec45d2f344ff01a9a922c))
* **docs:** add github social link to doc-site header ([37d0b59](https://github.com/Colin-XKL/FeedCraft/commit/37d0b59745aa867a466ecf48eba60c39f3a6b207))
* enhance LLM filters to use title and format payload in markdown codeblocks ([#654](https://github.com/Colin-XKL/FeedCraft/issues/654)) ([d191225](https://github.com/Colin-XKL/FeedCraft/commit/d1912252b8bd6e1641e59a3b5394731d7f7383d1))
* **feed-viewer:** enhance preview with HTML mode, card expand/collapse, and detail modal ([#804](https://github.com/Colin-XKL/FeedCraft/issues/804)) ([a14e886](https://github.com/Colin-XKL/FeedCraft/commit/a14e886f99e959873350a301042cd716ca5c3826))
* **fetcher:** configure http client with timeout instead of using DefaultClient ([6318cca](https://github.com/Colin-XKL/FeedCraft/commit/6318cca3c6b3cb33760dbf5845e451b446dec8c5))
* fix i18n for zh-TW ([c13f529](https://github.com/Colin-XKL/FeedCraft/commit/c13f529bf93f408a05db135088983f24502f7217))
* implement topic feed CRUD API and admin UI ([15da70a](https://github.com/Colin-XKL/FeedCraft/commit/15da70a816f41e0324974c69f2d0d52c91c26ddd))
* **inbox-ui:** optimize integration guide into tabs and add AI prompt generator ([a347b24](https://github.com/Colin-XKL/FeedCraft/commit/a347b2427ce449255f260f474a90927c91c58219))
* **inbox:** add gc stats and cleanup endpoints ([6e32030](https://github.com/Colin-XKL/FeedCraft/commit/6e32030f93d260ec84a36fc13a046ae3e3aaf596))
* **log:** allow custom LOG_LEVEL and demote language detection log to trace ([f69d105](https://github.com/Colin-XKL/FeedCraft/commit/f69d1054072b6eb31e20004ccd1bcb149f9ba83b))
* make embedding model configuration mandatory and remove defaults ([06db2ba](https://github.com/Colin-XKL/FeedCraft/commit/06db2bab8ee0970dc3a56ad71b8109ac24cf2192))
* merge Feed Compare into RSS Viewer page ([f8a51a0](https://github.com/Colin-XKL/FeedCraft/commit/f8a51a01844a55856b5e2b84a3c620b5dba7c61c))
* optimistic caching for topic sub-feeds with per-URI health tracking ([#805](https://github.com/Colin-XKL/FeedCraft/issues/805)) ([a930c53](https://github.com/Colin-XKL/FeedCraft/commit/a930c533cae83e690f2b326e50b8fdb3a7c5d97d))
* optimize concept visualizer page([#633](https://github.com/Colin-XKL/FeedCraft/issues/633)) ([a67876a](https://github.com/Colin-XKL/FeedCraft/commit/a67876a447b56a237468bc28d2d96e834099490f))
* optimize topic feed timestamp logic ([fb80f5d](https://github.com/Colin-XKL/FeedCraft/commit/fb80f5dd9466fc8fc8ad685185d42a1b44731dfa))
* prevent editing AtomCraft name after creation ([#801](https://github.com/Colin-XKL/FeedCraft/issues/801)) ([3d438a9](https://github.com/Colin-XKL/FeedCraft/commit/3d438a9c53c64ff822b1da417e198a39d7ee1606))
* refactor user-agent header logic and support custom ua via env var ([f8f31c7](https://github.com/Colin-XKL/FeedCraft/commit/f8f31c7ffd5f5cbdb6ba10c8cb2ebffeaaef2c03))
* rename curl-to-rss to json-to-rss to make it more clear ([5674f75](https://github.com/Colin-XKL/FeedCraft/commit/5674f75d8dfb36a2b1d77ab7a859729bb8f8c3f0))
* **router:** add public topic feed route ([3efcf01](https://github.com/Colin-XKL/FeedCraft/commit/3efcf01b3f16b4b931d570954306969ed04cb2fa))
* **source:** add web monitor source and preview tool ([65e797a](https://github.com/Colin-XKL/FeedCraft/commit/65e797a498ac4b3a75d68250f819ff4a6e20920c))
* support HTML to RSS navigation actions ([#833](https://github.com/Colin-XKL/FeedCraft/issues/833)) ([fbbc618](https://github.com/Colin-XKL/FeedCraft/commit/fbbc618c4361c39e4e3f951128e170f8158761a7))
* **topic:** 新增三种去重策略 —— 按标题、SimHash、Embedding ([#788](https://github.com/Colin-XKL/FeedCraft/issues/788)) ([ce77b15](https://github.com/Colin-XKL/FeedCraft/commit/ce77b15eff7d676f591c68c548990ca9f251c32d))
* unhide topic feed page + redesign input source editor ([#783](https://github.com/Colin-XKL/FeedCraft/issues/783)) ([e64ac9d](https://github.com/Colin-XKL/FeedCraft/commit/e64ac9dc1346b5d444e35f43008e826424bc3c4e))
* unhide topic feed page and re-enable detail navigation ([1d75166](https://github.com/Colin-XKL/FeedCraft/commit/1d7516602b17d3251876f40ef7502bb67c7723d8))
* unify feed preview targets ([#780](https://github.com/Colin-XKL/FeedCraft/issues/780)) ([66c9334](https://github.com/Colin-XKL/FeedCraft/commit/66c9334fc4af1e56c969748d0f728e1099336a33))
* update admin panel for topic feed management ([e8c24a8](https://github.com/Colin-XKL/FeedCraft/commit/e8c24a83e6628fa8134dab0d25789b1fee7cfc7c))
* update arch doc ([6273626](https://github.com/Colin-XKL/FeedCraft/commit/62736260654f6043ac554b4e354a7ac40aae0ffa))
* update default FC_LLM_MAX_CONCURRENCY to 3 ([4bd5992](https://github.com/Colin-XKL/FeedCraft/commit/4bd5992819f6a9701353dc274035ff6a7aec0bb1))
* update favicon support for admin panel ([df9ac9e](https://github.com/Colin-XKL/FeedCraft/commit/df9ac9e595332d311f9f8f0b5f185debec8a3a1a))
* update ui for custom recipe page ([826b6d1](https://github.com/Colin-XKL/FeedCraft/commit/826b6d1376766c1cc74d2e1ac3abc3117b8cd155))
* use FNV hash for content keys  to get bettert performance ([c7e60a4](https://github.com/Colin-XKL/FeedCraft/commit/c7e60a4b4b02599ac3cb23e68c1fb0f385008d08))
* **versioning:** release-please 流程与构建期版本展示 (COL-36) ([#909](https://github.com/Colin-XKL/FeedCraft/issues/909)) ([9de2c09](https://github.com/Colin-XKL/FeedCraft/commit/9de2c09019445b4c7fca6487d3e8945cb01f471e))
* 优化 Topic Feed 管理：分步向导、输入备注与预览 ([#806](https://github.com/Colin-XKL/FeedCraft/issues/806)) ([85838a2](https://github.com/Colin-XKL/FeedCraft/commit/85838a25118e5d23c48c5e9116a4e115ef5275c3))
* 将剩余 AI/Embedding craft 迁到原生 CraftFeed 处理器 ([#900](https://github.com/Colin-XKL/FeedCraft/issues/900)) ([2f4a245](https://github.com/Colin-XKL/FeedCraft/commit/2f4a245b4a6fa4fd3d32bbea0e454dd1fb22fd69))
* 限制 Feed 输入条目数 ([#881](https://github.com/Colin-XKL/FeedCraft/issues/881)) ([92ada8e](https://github.com/Colin-XKL/FeedCraft/commit/92ada8e2706063ac1fad52e01b483cf1dff96fe4))


### Bug Fixes

* add context timeout to feed fetch and retry loop ([bc4f3ac](https://github.com/Colin-XKL/FeedCraft/commit/bc4f3acb26ce90b83c0fcaf92cc8363a5c57140d))
* add generic MD5 hash helper ([83bdc46](https://github.com/Colin-XKL/FeedCraft/commit/83bdc46468bcf9fd3e2a35d18db154bc27cc2ad1))
* **admin:** HTML 转 RSS 允许抓取内网源并展示可读错误 ([0e5f2ea](https://github.com/Colin-XKL/FeedCraft/commit/0e5f2ea55e3bba168e7116545765ea17a696be28))
* **admin:** Topic 编辑页输入源与聚合规则改为纵向列表 ([597150a](https://github.com/Colin-XKL/FeedCraft/commit/597150a4148fdcbd6759c662ecf5f07a32a3fa2d))
* **admin:** 为空列表页增加创建引导 ([#880](https://github.com/Colin-XKL/FeedCraft/issues/880)) ([5c0b77e](https://github.com/Colin-XKL/FeedCraft/commit/5c0b77ef27f5af53a68d08794adc1559e5d39546))
* **admin:** 优化 Topic 向导、支持 Inbox，并修复非法 XML 字符解析失败 ([#876](https://github.com/Colin-XKL/FeedCraft/issues/876)) ([67e6370](https://github.com/Colin-XKL/FeedCraft/commit/67e63703496c5b116fbebabb03680ed00b6fe380))
* **admin:** 修复全部 Craft 列表名称列英文中途折行 ([3d99dbc](https://github.com/Colin-XKL/FeedCraft/commit/3d99dbc389e49284fdf740100bdfe5d90d666390))
* **admin:** 修复自定义配方表格单元格文本垂直堆叠 ([fa503a6](https://github.com/Colin-XKL/FeedCraft/commit/fa503a676c78df7b58636f7b4bccb720fa76f444))
* **admin:** 将欢迎页快速开始改为可交互 URL 生成 ([87238f3](https://github.com/Colin-XKL/FeedCraft/commit/87238f31325aa96832113f17b66df4c8e6577a96))
* **admin:** 居中 Topic 创建向导的核心内容区域 ([8d85111](https://github.com/Colin-XKL/FeedCraft/commit/8d85111373bcf11511c51aa2d2a0c29d8125b2ce))
* **admin:** 强制 Craft 名称列 nowrap，避免最长标识中途折行 ([d109fb6](https://github.com/Colin-XKL/FeedCraft/commit/d109fb6bf1fccfce2cad5445b2aa0308dae02b3f))
* **admin:** 文档中心菜单按当前 UI 语言跳转 (COL-37) ([#910](https://github.com/Colin-XKL/FeedCraft/issues/910)) ([245f21f](https://github.com/Colin-XKL/FeedCraft/commit/245f21f47e9f3365aefe18506749ecf461bbac14))
* **admin:** 自定义配方空表单保存前校验 (COL-27) ([#907](https://github.com/Colin-XKL/FeedCraft/issues/907)) ([c58693e](https://github.com/Colin-XKL/FeedCraft/commit/c58693ec2134426473ebe1e1a07261d0743700cd))
* **cache:** 用 singleflight 合并 CachedFunc 同 key 并发未命中 ([8ca02a5](https://github.com/Colin-XKL/FeedCraft/commit/8ca02a5b6ba17aacb7cd7469fc86c3f30c6a89d4))
* correctly handle logout action and provide feedback ([#741](https://github.com/Colin-XKL/FeedCraft/issues/741)) ([4dd29e2](https://github.com/Colin-XKL/FeedCraft/commit/4dd29e25d1889ab26a0ffbf8b1de48f5bee03ee8))
* correctly resolve relative feed links against base URL ([46c0f10](https://github.com/Colin-XKL/FeedCraft/commit/46c0f10ced127af47808ed6396641aaac30ac6e0))
* **craft:** 将上游抓取失败从 500 改为 502/504/422 ([2f27f3b](https://github.com/Colin-XKL/FeedCraft/commit/2f27f3bfb44de6bc9bfe61964283d05d2fea0217))
* **craft:** 浏览器访问 craft 端点时渲染可读预览 ([e0eaab9](https://github.com/Colin-XKL/FeedCraft/commit/e0eaab9065417326ab31bb3708a5ed00d05d7458))
* **craft:** 遗留适配器判空并统一 AI craft 构建期校验 ([#904](https://github.com/Colin-XKL/FeedCraft/issues/904)) ([bdd1c59](https://github.com/Colin-XKL/FeedCraft/commit/bdd1c593ee13830593335c6d0989d184f1076571))
* **craft:** 避免 fulltext-plus 因 browserless 超时全体失败 ([#883](https://github.com/Colin-XKL/FeedCraft/issues/883)) ([ab9e252](https://github.com/Colin-XKL/FeedCraft/commit/ab9e252aae3e5a2bce1f493e45940f00a58f537e))
* custom recipe and craft flow deletion behavior ([#715](https://github.com/Colin-XKL/FeedCraft/issues/715)) ([590f283](https://github.com/Colin-XKL/FeedCraft/commit/590f283ab8e28c9bb304d1e879442167527c27ac))
* **dao:** return nil on TopicFeed query failure ([9a0b797](https://github.com/Colin-XKL/FeedCraft/commit/9a0b7977aad40aaa7f7557394d2fd460d75c5d9f))
* **docs:** 修复 Inbox 指南页面显示异常 ([#886](https://github.com/Colin-XKL/FeedCraft/issues/886)) ([84cd9a1](https://github.com/Colin-XKL/FeedCraft/commit/84cd9a113382b34b00911ed8f29911b976a74ef3))
* enforce form validation before saving craft atom ([#642](https://github.com/Colin-XKL/FeedCraft/issues/642)) ([11227bc](https://github.com/Colin-XKL/FeedCraft/commit/11227bcc0688e0d378034763cc217f2232efb01c))
* Ensure valid JSON in jsonToRss placeholders ([93d5abf](https://github.com/Colin-XKL/FeedCraft/commit/93d5abf98cbb7c63011000b72e6ea64ad837a5f2))
* **html-to-rss:** use descendant combinator in generated CSS selectors ([1bc5c73](https://github.com/Colin-XKL/FeedCraft/commit/1bc5c73804d01ad824c1b39b41a9585f858e5aca))
* improve browserless error handling and feed viewer presentation ([#714](https://github.com/Colin-XKL/FeedCraft/issues/714)) ([154f911](https://github.com/Colin-XKL/FeedCraft/commit/154f911b99456a0290a6cf80b4dafd35e967fb3d))
* **model:** prevent nil article panics by filtering at source and adding checks ([53c26e0](https://github.com/Colin-XKL/FeedCraft/commit/53c26e075e2002197ebdd18e5bc24b298922f228))
* **model:** use fallback author from gofeed if authors list is empty ([40f54ad](https://github.com/Colin-XKL/FeedCraft/commit/40f54ad2e8e2a30ba57d0b3162362634e514a54e))
* normalize feed viewer errors ([c09136f](https://github.com/Colin-XKL/FeedCraft/commit/c09136f41b393fb7535682e46f1a2c7dfcf39866))
* normalize relative VITE_API_BASE_URL to start with a slash ([dfd822d](https://github.com/Colin-XKL/FeedCraft/commit/dfd822d4f936906553d0cb2b867e0c93f3962efa))
* **parser:** prevent clearing fallback content on HTML extraction failure ([6f38c1f](https://github.com/Colin-XKL/FeedCraft/commit/6f38c1f5978d15752c84803a477d756a2aaa7406))
* preserve feed viewer validation messages ([62ce2f0](https://github.com/Colin-XKL/FeedCraft/commit/62ce2f057c035df965f522ef53d6049d89e32a0d))
* **recipe:** add nil check for processedCraftFeed to prevent panic ([09d8bf1](https://github.com/Colin-XKL/FeedCraft/commit/09d8bf1fbc488b54d59f6e356fce5d307dfb22f7))
* **recipe:** 重复配方名称返回 409 而非暴露 SQL 错误 ([#888](https://github.com/Colin-XKL/FeedCraft/issues/888)) ([8c8a49e](https://github.com/Colin-XKL/FeedCraft/commit/8c8a49ee3900e8ade1fd60d9206f5795014e84c1))
* resolve inbox bugs found during testing ([d1c1262](https://github.com/Colin-XKL/FeedCraft/commit/d1c12626fe1a8249dc908d2e92adc87be5a9376e))
* restore description fallback before content-based crafts run ([c8a93c1](https://github.com/Colin-XKL/FeedCraft/commit/c8a93c1d88275897ad4922cbe3c8b763095ec4d6))
* return 404 on deleting non-existent topic feed ([e77c4bb](https://github.com/Colin-XKL/FeedCraft/commit/e77c4bb16b37ae9daad154fd393703d8aabd5f49))
* **util:** 修复 Markdown/HTML 转换中的段落换行丢失 ([#807](https://github.com/Colin-XKL/FeedCraft/issues/807)) ([e5ba366](https://github.com/Colin-XKL/FeedCraft/commit/e5ba366f4369d39db5529fd801154d10bc4d3cb9))
* **web-monitor:** 修复4处网页监控功能 Bug ([#789](https://github.com/Colin-XKL/FeedCraft/issues/789)) ([f8d95fb](https://github.com/Colin-XKL/FeedCraft/commit/f8d95fb52219951b10d3e2071a5aa08e6ce0231d))
* **web/admin:** handle copy link failure gracefully in TopicFeed detail view ([5a7e72a](https://github.com/Colin-XKL/FeedCraft/commit/5a7e72a140c664d3260df4204106fc37fd00b48a))


### Performance Improvements

* **craft:** 单 feed 文章并发进入全局 LLM 调度器 ([#902](https://github.com/Colin-XKL/FeedCraft/issues/902)) ([07d0aa8](https://github.com/Colin-XKL/FeedCraft/commit/07d0aa800b6b45a01b1249e90e3ac23e502e9cdc))

## [v3.2.0] (since v3.1.0)

### ⚠️ 破坏性变更 (Breaking Changes)

- Embedding 模型必须显式配置 `FC_EMBEDDING_API_MODEL`，不再提供默认值

### ✨ 新特性 (Features)

- **收件箱 (Inbox)**：脚本、自动化工具或第三方平台可以直接把文章 POST 给 FeedCraft，再作为 RSS 订阅。适合给本身没有 RSS 的来源（如 webhook 通知、私有数据）建一个订阅源；配套系统授权令牌、条目上限与过期清理
- **主题订阅 (Topic Feed)**：正式开放入口，把多个来源合成一个订阅。新增按标题、SimHash（字面近似）、Embedding（语义）三种跨源去重，不同媒体报道同一件事只看到一条；分步向导可为每个输入源加备注并即时预览效果，子源抓取失败时用缓存兜底而不是整个主题变空
- **新增原子工艺**：
  - `ai-filter`：用自然语言写筛选规则，让 LLM 决定文章留还是丢
  - `ai-content-process`：用自己的提示词加工正文，结果可插到文首、文末或替换原文
  - `re-title`：让 LLM 根据正文重写标题，改善标题党或含义不明的标题
  - `link-flatten`：把文章里的链接抽出来变成独立条目，适合订阅周报、导航页这类“链接列表”式内容
- **网页监控**：为没有 RSS 的网页建订阅，只在关注的字段（如价格、版本号）变化时才推送更新
- **RSS 生成增强**：HTML 转 RSS 支持抓取前先执行点击、翻页等导航动作，能取到需要交互才显示的内容，并可指定条目图标来源；JSON 转 RSS 支持 `$` 模板写法
- **预览体验**：RSS 预览支持 HTML 渲染模式、卡片展开/收起与详情弹窗，Feed 对比也并入同一页面，调 Craft 时可直接看到处理前后差异；浏览器服务改为可配置 Provider，不再限定 Browserless
- **输入条目上限**：可限制单次从源读取的条目数，订阅历史很长的大源时不必每次全量处理，明显减少耗时与 LLM 花费

### 🐛 问题修复 (Bug Fixes)

- `fulltext-plus`：个别文章渲染超时不再让整个 Feed 一起失败，浏览器取不到内容时自动回退普通全文提取
- 上游源抓取失败时返回更准确的状态码（502/504/422）；配方名重复时给出明确提示而非数据库错误
- 大量文章同时调用 LLM 时排队更平稳，正文里的 base64 图片会先清理以减少 token 消耗，缓存未命中时的重复请求会被合并
- 修复网页监控若干问题、Markdown/HTML 转换丢失段落换行、含非法 XML 字符的源解析失败、HTML 转 RSS 无法抓取内网源且错误提示不可读
- AtomCraft 创建后不再允许改名，避免引用它的 FlowCraft 和配方失效

### 📝 文档与杂项 (Documentation & Chores)

- Inbox / Topic / AI 内容处理文档（en / zh-CN / zh-TW）；开发环境指南
- 前端 TypeScript 工具链升级；CI 升级 Node 24，适配 Go 1.25 与 golangci-lint

## [v3.1.0] (since v3.0.0)

### ✨ 新特性 (Features & Refactors)

- **调度与并发优化 (Dispatcher & Concurrency)**:
  - 增加 LLM 并发控制逻辑 (concurrency control)
  - fulltext crafts 类新增 domain 级别 rate limiting
- **搜索与 RSS (Search & RSS)**:
  - Search-to-RSS: 增加 enhanced mode
  - RSS wizards: 使用 `limax` 自动生成 recipe ID
- **系统与用户体验**:
  - 增加 System Health Check 页面
  - 改进错误提示用户体验 (better error ux)

### 🐛 问题修复 (Bug Fixes)

- 修复 `priority dispatcher` 中的潜在死锁问题
- 修复 search provider configs 中的 jq expressions 错误
- 允许清除 search provider API key
- settings: 修复读取现有 search provider config 时的错误处理
- 修复 search provider active check 并增加 timeout
- monitor: 修复并暴露 search provider check 中的 db errors

### 📝 文档与杂项 (Documentation & Chores)

- **文档 (Docs)**:
  - 增加系统工具的文档 (viewer, compare, health)
  - 增加比较 FeedCraft 与其他工具的文档
  - 增加繁体中文文档 (zh-tw doc)
  - 更新 minimal docker compose example, theme 和 badge

## [v3.0.0] (since v2.1)

### ⚠️ 破坏性变更 (Breaking Changes)

- **LLM 配置更新**: 重构了 LLM 集成，引入了通用的环境变量配置
  - 新增 `FC_LLM_API_TYPE` (支持 `openai`, `ollama`)
  - `FC_OPENAI_ENDPOINT` 重命名为 `FC_LLM_API_BASE`
  - `FC_OPENAI_AUTH_KEY` 重命名为 `FC_LLM_API_KEY`
  - `FC_OPENAI_DEFAULT_MODEL` 重命名为 `FC_LLM_API_MODEL`
  - 旧变量仍暂时兼容但有废弃警告
- **术语变更**: UI 和文档中的 "Craft Atom" 重命名为 "AtomCraft" (原子工艺), "Craft Flow" 重命名为 "FlowCraft" (组合工艺)

### ✨ 新特性 (Features)

- **RSS 生成器工具集**:
  - 新增 **HTML 转 RSS** 工具: 支持交互式选择器拾取、增强模式 (无头浏览器/Browserless)、富文本预览及智能选择逻辑
  - 新增 **JSON 转 RSS** 工具: 支持通过 JQ 表达式从 JSON 源生成 RSS
  - 新增 **搜索 转 RSS** 工具: 集成 SearXNG 和 LiteLLM，支持通过搜索结果生成 RSS
  - 新增 **快速开始 (URL 生成器)**: 支持生成和解析 FeedCraft URL
  - Curl 转 RSS 支持配置 HTTP 方法和请求体
- **工艺组件 (Atom/Flow)**:
  - 新增 `time-limit` (时间限制) 原子工艺
  - 新增 `beautify-content` (内容美化) 原子工艺
  - 新增 `article-summary` (文章摘要) 原子工艺
  - 新增 `immersive-translate` (沉浸式翻译) 组合工艺
  - 新增通用 `llm-filter` (LLM 过滤器) 原子工艺
  - `fulltext-plus`: 增加 `wait` (等待时间) 和 `mode` (如 `networkidle2`) 参数以更好支持动态网页
  - 支持 `DEFAULT_TARGET_LANG` 环境变量，用于控制翻译目标语言
- **用户界面与体验 (UI/UX)**:
  - 新增 **服务依赖状态** 页面，用于监控 SQLite, Redis, Browserless, LLM 等服务状态
  - 新增 **搜索提供商设置** 页面
  - 应用自定义 Arco Design 主题
  - 重构 **Craft Flow 编辑器**: 采用列表式编辑，支持拖拽排序
  - 改进 **Craft 选择器**: 模块化拆分，支持分类展示和多选
  - 自定义配方编辑器: 支持 JSON 格式化、一键复制配置
  - 添加关键操作的确认对话框 (如删除)
- **基础设施与后端**:
  - 支持 Ollama 作为 LLM 提供商
  - 支持配置多个 LLM 模型并实现自动重试逻辑
  - 优化 LLM 调用: 增加内容处理选项 (移除链接/图片) 以节省 Token
  - 构建流程: 注入版本、提交哈希等元数据到二进制文件
  - 新增 GitHub Actions CI 工作流

### 🐛 问题修复 (Bug Fixes)

- **HTML 转 RSS**:
  - 修复空响应导致的静默失败，优化错误处理
  - 修复向导中 Axios 响应未正确解包的问题
  - 优化 Fetch 逻辑，增加 User-Agent 和标准头以减少被拦截概率
- **搜索转 RSS**:
  - 修复生成失败时返回 200 状态码的问题 (现返回 500)
  - 处理数据库读取配置失败的情况
- **系统与路由**:
  - 修复缺失的 API 路由返回 HTML 的问题 (现返回 404 JSON)
  - 修复无效内存地址引用导致的 Panic
  - 验证 Browserless 服务返回的 HTTP 状态码
- **其他**:
  - 修复 RSS 生成器 CSS 预览问题
  - 修复 Docker 发布工作流中的 helper 错误

### 📝 文档与杂项 (Documentation & Chores)

- **文档**:
  - 新增关于搜索转 RSS、JSON 转 RSS、系统原子工艺的详细指南
  - 更新快速开始和自定义配置文档
  - 重构文档结构，迁移至 Astro Starlight
- **依赖与构建**:
  - 升级 Web 端 Vite 至 v5, TypeScript 至 v5
  - 升级 Go 和 Node.js 依赖 (如 gorm, axios, vue-router 等)
  - 更新 `.gitignore` 和 `Taskfile`
