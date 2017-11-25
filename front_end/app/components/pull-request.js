import Ember from 'ember';
const {Component, computed} = Ember;

import PropTypeMixin, {PropTypes} from 'ember-prop-types';
import Pull from '../models/pull';

export default Component.extend(PropTypeMixin, {
  propTypes: {
    pull: PropTypes.instanceOf(Pull),
  },

  comments: computed.alias('pull.comments'),
})