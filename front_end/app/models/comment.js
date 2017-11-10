import DS from 'ember-data';

export default DS.Model.extend({
  body: DS.attr('string'),
  position: DS.attr('number'),
  originalPosition: DS.attr('number'),
  user: DS.belongsTo('user'),
})