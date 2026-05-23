export default {
  // Inbox Management
  'inbox.title': '收件箱管理',
  'inbox.desc':
    '管理您的第三方主動資料投遞收件箱。收件箱接受標準 JSON 資料格式投遞並儲存。',
  'inbox.btn.create': '新建收件箱',
  'inbox.btn.edit': '編輯收件箱',
  'inbox.btn.delete': '刪除',
  'inbox.btn.guide': '整合指南',
  'inbox.id': '收件箱 ID',
  'inbox.id.placeholder': '請輸入唯一的收件箱 ID (例如: test-inbox)',
  'inbox.name': '標題',
  'inbox.name.placeholder': '請輸入收件箱標題',
  'inbox.description': '描述',
  'inbox.description.placeholder': '可選的描述資訊',
  'inbox.maxItems': '最大保存數',
  'inbox.isPublic': '公開可見性',
  'inbox.isPublic.true': '公開 (無需鑑權查看文章)',
  'inbox.isPublic.false': '私有 (必須提供 Token 查看)',
  'inbox.status.public': '公開',
  'inbox.status.private': '私有',
  'inbox.actions': '操作',
  'inbox.deleteConfirm':
    '確定要刪除此收件箱嗎？該收件箱內的所有文章也將被連帶刪除！',
  'inbox.deleteSuccess': '收件箱刪除成功',
  'inbox.createSuccess': '收件箱建立成功',
  'inbox.updateSuccess': '收件箱更新成功',

  // Integration Guide Modal
  'inbox.guide.title': '第三方整合與資料推送指南',
  'inbox.guide.pushUrl': '推送 API 位址',
  'inbox.guide.pushUrl.desc':
    '僅支援標準 HTTP POST 請求投遞資料。必須包含 SystemAuthToken 鑑權標頭。',
  'inbox.guide.headers': 'HTTP 請求標頭 (Headers)',
  'inbox.guide.body': '請求體 JSON 格式 (Array)',
  'inbox.guide.example': 'cURL 推送範例',
  'inbox.guide.query': 'RSS 訂閱與正文拉取',
  'inbox.guide.query.public': '該收件箱為公開狀態，任何人皆可直接拉取或閱讀。',
  'inbox.guide.query.private':
    '該收件箱為私有狀態，訂閱或拉取文章正文必須在 URL Query 帶上 `?token=YOUR_TOKEN` 進行身份驗證。',

  // Recipe integration guide steps
  'inbox.guide.recipe.heading': '透過配方訂閱 RSS',
  'inbox.guide.recipe.step1.pre': '前往',
  'inbox.guide.recipe.step1.link': '「自訂配方」',
  'inbox.guide.recipe.step1.post': '頁面，建立一個新配方。',
  'inbox.guide.recipe.step2': '將資料來源類型 (Source Type) 設為',
  'inbox.guide.recipe.step3':
    '在 Source Config JSON 中，設定當前收件箱的 ID 對應：',
  'inbox.guide.recipe.step4.pre': '在配方列表點擊',
  'inbox.guide.recipe.step4.link': '『複製連結』',
  'inbox.guide.recipe.step4.post': '即可獲取完整的 RSS 訂閱位址！',

  // System Auth Token Management
  'systemAuthToken.title': '系統授權令牌',
  'systemAuthToken.desc':
    '產生並管理全域的系統授權令牌。支援在不同外部裝置、自動化腳本或第三方平台上獨立使用。',
  'systemAuthToken.btn.create': '產生新令牌',
  'systemAuthToken.id': 'ID',
  'systemAuthToken.label': '令牌標籤 (用途說明)',
  'systemAuthToken.label.placeholder':
    '請輸入令牌的用途 (例如: iPhone 推送, HomeAssistant)',
  'systemAuthToken.value': 'Token 密鑰',
  'systemAuthToken.createdAt': '建立時間',
  'systemAuthToken.actions': '操作',
  'systemAuthToken.deleteConfirm':
    '確定要廢除並刪除此授權令牌嗎？刪除後，使用該令牌的外部整合推送或讀取將立即失效！',
  'systemAuthToken.deleteSuccess': '授權令牌廢除成功',
  'systemAuthToken.createSuccess': '新令牌產生成功！',
  'systemAuthToken.createdAlert.title': '請安全備份您的新授權令牌',
  'systemAuthToken.createdAlert.desc':
    '這是該令牌唯一一次以明文形式展示。為了您的系統安全，關閉此視窗後將無法再次找回明文，請立即複製並妥善保管！',
  'systemAuthToken.copied': '已成功複製到剪貼簿',
  'systemAuthToken.btn.copy': '複製',
  'systemAuthToken.btn.close': '關閉',
};
