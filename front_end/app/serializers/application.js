import DS from 'ember-data';
import { singularize } from 'ember-inflector';

export default DS.RESTSerializer.extend({
  normalizeQueryRecordResponse(store, primaryModelClass, payload) {
    for (const key of Object.keys(payload)) {
      if (this.modelNameFromPayloadKey(key) === primaryModelClass.modelName && Array.isArray(payload[key])) {
        payload[singularize(key)] = payload[key][0];
        delete payload[key];
      }
    }
    return this._super(...arguments);
  },
});