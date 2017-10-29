import Ember from 'ember';
import config from './config/environment';

const Router = Ember.Router.extend({
  location: config.locationType,
  rootURL: config.rootURL
});

Router.map(function() {
  this.route('repos', { path: '/' });
  this.route('repo', { path: '/repo/:repo_id' });
  this.route('pulls');
  this.route('pull', { path: '/pull/:pull_id' });
});

export default Router;
