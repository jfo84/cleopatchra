import DS from 'ember-data';

export default DS.Model.extend({
  url: DS.attr('string'),
  body: DS.attr('string'),
  state: DS.attr('string'),
  title: DS.attr('string'),
})