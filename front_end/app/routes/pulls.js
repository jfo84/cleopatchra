import Ember from 'ember';

export default Ember.Route.extend({
  async model() {
    try {
      const pulls = await this.get('store').query('pull', { page: 1, limit: 10 });
      return { ...pulls };
    } catch (e) {
      console.log(e);
    }
  },
});
