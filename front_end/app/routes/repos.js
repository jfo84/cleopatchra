import Ember from 'ember';

export default Ember.Route.extend({
  async model() {
    try {
      const repos = await this.get('store').query('repo', { page: 1, limit: 10 });
      return { repos };
    } catch (e) {
      console.log(e);
    }
  },
});