export default {
  // Inbox Management
  'inbox.title': 'Inbox Management',
  'inbox.desc':
    'Manage your third-party data-push inboxes. Inboxes accept standard JSON arrays to store articles.',
  'inbox.btn.create': 'Create Inbox',
  'inbox.btn.edit': 'Edit Inbox',
  'inbox.btn.delete': 'Delete',
  'inbox.btn.guide': 'Integration Guide',
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
};
