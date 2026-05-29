export default {
  // Inbox Management
  'inbox.title': 'Inbox Management',
  'inbox.desc':
    'Manage your third-party data-push inboxes. Inboxes accept standard JSON arrays to store articles.',
  'inbox.btn.create': 'Create Inbox',
  'inbox.btn.edit': 'Edit Inbox',
  'inbox.btn.delete': 'Delete',
  'inbox.btn.guide': 'Integration Guide',
  'inbox.btn.preview': 'Preview',
  'inbox.id': 'Inbox ID',
  'inbox.id.placeholder': 'Enter unique Inbox ID (e.g., test-inbox)',
  'inbox.name': 'Title',
  'inbox.name.placeholder': 'Enter inbox title',
  'inbox.description': 'Description',
  'inbox.description.placeholder': 'Optional description',
  'inbox.maxItems': 'Max Items Limit',
  'inbox.isPublic': 'Public Access',
  'inbox.isPublic.true': 'Public (anonymous reads allowed)',
  'inbox.isPublic.false': 'Private (auth token required for reading)',
  'inbox.status.public': 'Public',
  'inbox.status.private': 'Private',
  'inbox.actions': 'Actions',
  'inbox.deleteConfirm':
    'Are you sure you want to delete this inbox? All stored articles will be permanently removed as well!',
  'inbox.deleteSuccess': 'Inbox deleted successfully',
  'inbox.createSuccess': 'Inbox created successfully',
  'inbox.updateSuccess': 'Inbox updated successfully',

  // Integration Guide Modal
  'inbox.guide.title': 'Integration & Push Guide',
  'inbox.guide.pushUrl': 'Push API URL',
  'inbox.guide.pushUrl.desc':
    'Supports standard HTTP POST requests only. SystemAuthToken header is mandatory.',
  'inbox.guide.headers': 'HTTP Headers',
  'inbox.guide.body': 'Request Body JSON Format (Array)',
  'inbox.guide.example': 'cURL Push Code Example',
  'inbox.guide.query': 'RSS Subscription & Reading',
  'inbox.guide.query.public':
    'This inbox is Public. Anyone can fetch or subscribe to it directly.',
  'inbox.guide.query.private':
    'This inbox is Private. You must attach `?token=YOUR_TOKEN` query parameter to subscribe or view articles.',

  // Integration Guide Modal Tabs & Prompt Generator
  'inbox.guide.tab.push': 'How to Push to Inbox',
  'inbox.guide.tab.subscribe': 'How to Subscribe to RSS',
  'inbox.guide.prompt.title': 'AI Agent Script Prompt Generator',
  'inbox.guide.prompt.desc':
    'Generate a specialized System Prompt to guide external AI Agents (e.g., ChatGPT, Claude) in writing web scraping and data push scripts.',
  'inbox.guide.prompt.btn': 'Generate Prompt',
  'inbox.guide.prompt.preview': 'Prompt Preview',
  'inbox.guide.prompt.copy': 'Copy Prompt',
  'inbox.guide.prompt.template': `You are an expert AI assistant specializing in web data scraping. Your task is to write an automation script (e.g., in Python or Node.js) to scrape data from a specific website or data source and push it directly to the feed-craft Inbox.

[API Information]
- Push API URL: {pushUrl}
- Authorization Header: Bearer <YOUR_SYSTEM_AUTH_TOKEN>
- Request Method: POST
- Content-Type: application/json
- Request Body Format (JSON Array):
{jsonSample}

[Requirements]
1. Write a script to scrape data from the target website (you may use requests/BeautifulSoup/Playwright for Python, or axios/puppeteer for Node.js).
2. Format and organize the scraped data into the JSON array structure specified above.
3. Call the Push API URL, attaching the Authorization header, and send the data using HTTP POST to feed-craft.
4. Include robust error handling and logging to ensure the data is pushed successfully.
5. The script should support scheduled runs or single-run invocation.

Please write an efficient, robust scraping and pushing script based on the requirements above.`,

  // Direct RSS URL
  'inbox.guide.directRss.desc':
    'Subscribe directly with this RSS URL — no Custom Recipe needed:',
  'inbox.guide.directRss.privateHint':
    'Private inbox — append token to the URL:',

  // Recipe integration guide steps
  'inbox.guide.recipe.heading': 'RSS Subscription via Recipe',
  'inbox.guide.recipe.step1.pre': 'Go to the',
  'inbox.guide.recipe.step1.link': 'Custom Recipe',
  'inbox.guide.recipe.step1.post': 'page and create a new Recipe.',
  'inbox.guide.recipe.step2': 'Set the Source Type to',
  'inbox.guide.recipe.step3':
    'In the Source Config JSON field, configure the inbox ID mapping:',
  'inbox.guide.recipe.step4.pre': 'Click',
  'inbox.guide.recipe.step4.link': 'Copy Link',
  'inbox.guide.recipe.step4.post':
    'in the Recipe list to get the RSS feed URL!',

  // System Auth Token Management
  'systemAuthToken.title': 'System Auth Token',
  'systemAuthToken.desc':
    'Generate and manage global system authentication tokens. They can be used independently across different devices, scripts, or integrations.',
  'systemAuthToken.btn.create': 'Generate Token',
  'systemAuthToken.id': 'ID',
  'systemAuthToken.label': 'Token Label (Purpose)',
  'systemAuthToken.label.placeholder':
    'Enter token purpose (e.g., iPhone Push, HomeAssistant Monitor)',
  'systemAuthToken.value': 'Token Key',
  'systemAuthToken.createdAt': 'Created At',
  'systemAuthToken.actions': 'Actions',
  'systemAuthToken.deleteConfirm':
    'Are you sure you want to delete and revoke this token? All automated integrations using this token will stop working immediately!',
  'systemAuthToken.deleteSuccess': 'Token revoked successfully',
  'systemAuthToken.createSuccess': 'Token generated successfully!',
  'systemAuthToken.createdAlert.title':
    'Please Securely Back Up Your New Auth Token',
  'systemAuthToken.createdAlert.desc':
    'This is the only time your token is displayed in plain text. For safety reasons, you cannot retrieve it after closing this window. Please copy it immediately!',
  'systemAuthToken.copied': 'Copied to clipboard successfully',
  'systemAuthToken.btn.copy': 'Copy',
  'systemAuthToken.btn.close': 'Close',
};
