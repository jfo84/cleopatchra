import Ember from 'ember';

const { Controller, computed } = Ember;

export default Controller.extend({
  pull: computed.alias('model.pull'),
});