import DS from 'ember-data';
import { singularize } from 'ember-inflector';

export default DS.JSONAPISerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    /**
     * This hack keeps the server from having to store anything except the raw GitHub API response underlying each object.
     * Without it, the seeder would have to add a type to each response, which means un-marshalling and marshalling every pass.
     * Maybe it doesn't matter since the seeder already hits the API request limit even though ruby is slow as shit.
     * 
     * In the future (hopefully); with things like a condensed tree format for a PR, search indexing, and an agnostic
     * approach across git providers; a solution would be needed for this.
     */
    const type = Object.keys(payload)[0];
    const modelArray = payload[type];
    const typedModels = modelArray.map(model => {
      // It seems like a better idea to send a route-specific key (i.e. repos, pulls instead of repo, pull)
      // The information could be useful later
      model['type'] = singularize(type);

      return model;
    });
    const normalizedPayload = { data: typedModels };

    return this._super(store, primaryModelClass, normalizedPayload, id, requestType);
  },
});