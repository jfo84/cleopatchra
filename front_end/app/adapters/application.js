import DS from 'ember-data';
import UrlTemplates from 'ember-data-url-templates';

export default DS.JSONAPIAdapter.extend(UrlTemplates, {
  host: 'http://localhost:7000',
});