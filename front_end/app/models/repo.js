import DS from 'ember-data';

export default DS.Model.extend({
  pulls: DS.hasMany('pull'),
  url: DS.attr('string'),
  name: DS.attr('string'),
})