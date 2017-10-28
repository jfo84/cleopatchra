import DS from 'ember-data';

export default DS.Model.extend({
  id: DS.attr('number'),
  url: DS.attr('string'),
  body: DS.attr('string'),
  state: DS.attr('string'),
  title: DS.attr('string'),
})