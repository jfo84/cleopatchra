import Ember from 'ember';

const { Controller, computed } = Ember;

export default Controller.extend({
  pulls: computed.alias('model.pulls'),
});
