import Ember from 'ember';

export default Ember.Route.extend({
  async model(params) {
    try {
      const pull = await this.get('store').findRecord('pull', params['pull_id']);
      return { pull };
    } catch (e) {
      console.log(e);
    }
  },
});