import DS from 'ember-data';

export default DS.Model.extend({
  repo: DS.belongsTo('repo'),
  comments: DS.hasMany('comment'),
  url: DS.attr('string'),
  body: DS.attr('string'),
  state: DS.attr('string'),
  title: DS.attr('string'),
  number: DS.attr('number'),
})