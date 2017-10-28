import Ember from 'ember';

const { Controller, computed } = Ember;

export default Controller.extend({
  repos: computed.alias('model.pulls'),
});
