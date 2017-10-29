import DS from 'ember-data';

export default DS.Model.extend({
  url: DS.attr('string'),
  name: DS.attr('string'),
  organization: DS.attr(),
})