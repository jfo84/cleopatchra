import Ember from 'ember';
const {Component} = Ember;

import PropTypeMixin, {PropTypes} from 'ember-prop-types';
import Pull from '../models/pull';
import Comment from '../models/comment';

export default Component.extend(PropTypeMixin, {
  propTypes: {
    pull: PropTypes.instanceOf(Pull),
    comments: PropTypes.arrayOf(PropTypes.instanceOf(Comment)),
  },
})