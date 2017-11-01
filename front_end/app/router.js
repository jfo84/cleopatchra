import Ember from 'ember';
import config from './config/environment';

const Router = Ember.Router.extend({
  location: config.locationType,
  rootURL: config.rootURL
});

Router.map(function() {
  this.route('repos', { path: '/' });
  this.route('repos', { path: '/repos/:repo_id' }, function() {
    this.route('pulls');
  });
  this.route('pulls', { path: '/pulls/:pull_id' });
});

export default Router;
