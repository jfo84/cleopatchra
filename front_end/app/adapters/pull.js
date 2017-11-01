import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  queryUrlTemplate: '{+host}/repos/{repoId}/pulls',
});
