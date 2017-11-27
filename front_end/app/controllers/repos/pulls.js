import Ember from 'ember';

const { Controller, computed } = Ember;

export default Controller.extend({
  rows: computed.alias('model.pulls'),
});
