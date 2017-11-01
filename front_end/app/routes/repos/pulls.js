import Ember from 'ember';

export default Ember.Route.extend({
  async model(params, transition) {
    const repoId = transition.params['repos']['repo_id'];
    try {
      const pulls = await this.get('store').query('pull', { page: 1, limit: 10, repoId: repoId });
      return { pulls };
    } catch (e) {
      console.log(e);
    }
  },
});
