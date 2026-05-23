export default {
  // Inbox Management
  'inbox.title': '收件箱管理',
  'inbox.desc':
    '管理您的第三方主动数据投递收件箱。收件箱接受标准 JSON 数据格式投递并存储。',
  'inbox.btn.create': '新建收件箱',
  'inbox.btn.edit': '编辑收件箱',
  'inbox.btn.delete': '删除',
  'inbox.btn.guide': '集成指南',
  'inbox.id': '收件箱 ID',
  'inbox.id.placeholder': '请输入唯一的收件箱 ID (例如: test-inbox)',
  'inbox.name': '标题',
  'inbox.name.placeholder': '请输入收件箱标题',
  'inbox.description': '描述',
  'inbox.description.placeholder': '可选的描述信息',
  'inbox.maxItems': '最大保存数',
  'inbox.isPublic': '公开可见性',
  'inbox.isPublic.true': '公开 (无需鉴权查看文章)',
  'inbox.isPublic.false': '私有 (必须提供 Token 查看)',
  'inbox.status.public': '公开',
  'inbox.status.private': '私有',
  'inbox.actions': '操作',
  'inbox.deleteConfirm':
    '确定要删除此收件箱吗？该收件箱内的所有文章也将被级联删除！',
  'inbox.deleteSuccess': '收件箱删除成功',
  'inbox.createSuccess': '收件箱创建成功',
  'inbox.updateSuccess': '收件箱更新成功',

  // Integration Guide Modal
  'inbox.guide.title': '第三方集成与数据推送指南',
  'inbox.guide.pushUrl': '推送 API 地址',
  'inbox.guide.pushUrl.desc':
    '仅支持标准 HTTP POST 请求投递数据。必须包含 SystemAuthToken 鉴权头。',
  'inbox.guide.headers': 'HTTP 请求头 (Headers)',
  'inbox.guide.body': '请求体 JSON 格式 (Array)',
  'inbox.guide.example': 'cURL 推送示例',
  'inbox.guide.query': 'RSS 订阅与正文拉取',
  'inbox.guide.query.public': '该收件箱为公开状态，任何人皆可直接拉取或阅读。',
  'inbox.guide.query.private':
    '该收件箱为私有状态，订阅或拉取文章正文必须在 URL Query 带上 `?token=YOUR_TOKEN` 进行身份验证。',

  // Recipe integration guide steps
  'inbox.guide.recipe.heading': '通过配方订阅 RSS',
  'inbox.guide.recipe.step1.pre': '前往',
  'inbox.guide.recipe.step1.link': '「自定义配方」',
  'inbox.guide.recipe.step1.post': '页面，创建一个新配方。',
  'inbox.guide.recipe.step2': '将数据源类型 (Source Type) 设为',
  'inbox.guide.recipe.step3':
    '在 Source Config JSON 中，配置当前收件箱的 ID 映射：',
  'inbox.guide.recipe.step4.pre': '在配方列表点击',
  'inbox.guide.recipe.step4.link': '『复制链接』',
  'inbox.guide.recipe.step4.post': '即可获取完整的 RSS 订阅地址！',

  // System Auth Token Management
  'systemAuthToken.title': '系统授权令牌',
  'systemAuthToken.desc':
    '生成并管理全局的系统授权令牌。支持在不同外部设备、自动化脚本、或者第三方平台上独立使用。',
  'systemAuthToken.btn.create': '生成新令牌',
  'systemAuthToken.id': 'ID',
  'systemAuthToken.label': '令牌标签 (用途说明)',
  'systemAuthToken.label.placeholder':
    '请输入令牌的用途 (例如: iPhone 推送, HomeAssistant)',
  'systemAuthToken.value': 'Token 密钥',
  'systemAuthToken.createdAt': '创建时间',
  'systemAuthToken.actions': '操作',
  'systemAuthToken.deleteConfirm':
    '确定要废除并删除此授权令牌吗？删除后，使用该令牌的外部集成推送或读取将会立即失效！',
  'systemAuthToken.deleteSuccess': '授权令牌废除成功',
  'systemAuthToken.createSuccess': '新令牌生成成功！',
  'systemAuthToken.createdAlert.title': '请安全备份您的新授权令牌',
  'systemAuthToken.createdAlert.desc':
    '这是该令牌唯一一次以明文形式展示。为了您的系统安全，关闭此窗口后将无法再次找回明文，请立即复制并妥善保管！',
  'systemAuthToken.copied': '已成功复制到剪贴板',
  'systemAuthToken.btn.copy': '复制',
  'systemAuthToken.btn.close': '关闭',
};
