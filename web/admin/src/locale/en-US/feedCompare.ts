export default {
  feedCompare: {
    description:
      'Specify the feed source URL and the craft (craft flow or atom) to be applied, and compare the results before and after processing.',
    inputLink: 'Input Link',
    placeholder: 'Enter RSS feed URL',
    compare: 'Compare',
    originalFeed: 'Original Feed',
    craftAppliedFeed: 'Craft Applied Feed',
    selectCraftFlow: {
      placeholder: 'Select CraftFlow',
      tabs: {
        system: 'System Craft Templates',
        user: 'User Craft Atoms',
        flow: 'Craft Flows',
      },
    },
    message: {
      inputRequired: 'Please enter a feed URL and select a craft',
      unknownError: 'Unknown Error',
    },
  },
};
