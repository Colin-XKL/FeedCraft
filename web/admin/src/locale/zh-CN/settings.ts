export default {
  'settings.title': '页面配置',
  'settings.themeColor': '主题色',
  'settings.content': '内容区域',
  'settings.search': '搜索',
  'settings.language': '语言',
  'settings.navbar': '导航栏',
  'settings.menuWidth': '菜单宽度 (px)',
  'settings.navbar.theme.toLight': '点击切换为亮色模式',
  'settings.navbar.theme.toDark': '点击切换为暗黑模式',
  'settings.navbar.screen.toFull': '点击切换全屏模式',
  'settings.navbar.screen.toExit': '点击退出全屏模式',
  'settings.navbar.alerts': '消息通知',
  'settings.menu': '菜单栏',
  'settings.topMenu': '顶部菜单栏',
  'settings.tabBar': '多页签',
  'settings.footer': '底部',
  'settings.otherSettings': '其他设置',
  'settings.colorWeak': '色弱模式',
  'settings.alertContent':
    '配置之后仅是临时生效，要想真正作用于项目，点击下方的 "复制配置" 按钮，将配置替换到 settings.json 中即可。',
  'settings.copySettings': '复制配置',
  'settings.copySettings.message':
    '复制成功，请粘贴到 src/settings.json 文件中',
  'settings.close': '关闭',
  'settings.color.tooltip':
    '根据主题颜色生成的 10 个梯度色（将配置复制到项目中，主题色才能对亮色 / 暗黑模式同时生效）',
  'settings.menuFromServer': '菜单来源于后台',
  'settings.searchProvider.title': '搜索提供商配置',
  'settings.searchProvider.alert':
    '配置“搜索转 RSS”功能使用的搜索提供商。支持 LiteLLM (代理) 和 SearXNG。',
  'settings.searchProvider.provider': '提供商实现',
  'settings.searchProvider.provider.tooltip':
    "要使用的内部提供商逻辑。支持 'litellm' 和 'searxng'。",
  'settings.searchProvider.apiUrl': 'API URL',
  'settings.searchProvider.apiUrl.tooltip':
    '搜索提供商的基准 URL (例如 http://litellm-proxy:4000 或 http://searxng.local)。',
  'settings.searchProvider.apiKey': 'API Key',
  'settings.searchProvider.apiKey.tooltip':
    '如果提供商需要身份验证，则为可选 API 密钥。',
  'settings.searchProvider.toolName': '工具名称 / 引擎',
  'settings.searchProvider.toolName.tooltip':
    "对于 LiteLLM：具体工具名称 (如 'google-search')。对于 SearXNG：逗号分隔的引擎列表 (如 'google,bing')。留空使用默认值。",
  'settings.searchProvider.save': '保存配置',
  'settings.searchProvider.configure': '配置搜索提供商',
  'settings.searchProvider.msg.urlRequired': 'API URL 是必填项',
  'settings.searchProvider.msg.saved': '配置已保存',
  'settings.searchProvider.msg.failed': '保存配置失败',
  'settings.searchProvider.placeholder.selectProvider': '选择提供商',
  'settings.searchProvider.placeholder.apiUrl': 'http://litellm-proxy:4000',
  'settings.searchProvider.placeholder.apiKey': 'sk-...',
  'settings.searchProvider.placeholder.toolName': 'google-search',
};
