import DS from 'ember-data';

export default DS.Model.extend({
  id: DS.attr('number'),
  url: DS.attr('string'),
  name: DS.attr('string'),
  organization: DS.attr(),
})