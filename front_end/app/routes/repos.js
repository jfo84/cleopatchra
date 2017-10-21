import Ember from 'ember';

export default Ember.Route.extend({
  async model(params) {
    try {
      const repos = await this.get('store').query('repos', { page: 1 });
      return { ...repos };
    } catch (e) {
      console.log(e);
    }
  },
});