import * as actionTypes from './actionTypes';
import queryString from 'query-string';
import dotenv from 'dotenv';

dotenv.config();
const authToken = process.env.API_TOKEN
const BASE_PARAMS = { authToken };

export const initialize = () => {
  return (dispatch) => {
    dispatch(fetchPulls());
  };
};

export const updatePulls = (repoId, page, limit) => {
  return (dispatch) => {
    dispatch(fetchPulls(repoId, page, limit))
  }
}

const requestPulls = () => {
  return {
    type: actionTypes.REQUEST_PULLS,
    payload: {
      isFetching: true
    }
  };
};

const receivePulls = (pulls) => {
  return {
    type: actionTypes.RECEIVE_PULLS,
    payload: {
      isFetching: false,
      pulls
    }
  };
};

const fetchPulls = (repoId, page = 1, limit = 10) => {
  return (dispatch) => {
    dispatch(requestPulls());

    const params = Object.assign(BASE_PARAMS, { page, limit });
    // TODO: Figure out passing repoId from the router
    var url = `http://localhost:7000/repos/10270250/pulls`;

    return fetch(url).then((response) => {
      return response.json();
    }).then((pulls) => {
      dispatch(receivePulls(pulls));
    });
  }
};

export const pageForwards = () => {
  return (dispatch, getState) => {
    var { page, limit } = getState();
    var page = page + limit;
    dispatch({
      type: actionTypes.CHANGE_PAGE,
      payload: page
    });
  }
};

export const pageBackwards = () => {
  return (dispatch, getState) => {
    var { page, limit } = getState();
    var page = page - limit;
    dispatch({
      type: actionTypes.CHANGE_PAGE,
      payload: page
    });
  };
};
