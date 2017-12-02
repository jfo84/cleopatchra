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

export const updatePulls = (repoId, page, numPulls) => {
  return (dispatch) => {
    dispatch(fetchPulls(repoId, page, numPulls))
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

    const params = Object.assign(BASE_PARAMS, { page, numPulls });
    var url = `localhost:7000/repos/${repoId}/pulls`;

    return fetch(url).then((simpleResponse) => {
      return simpleResponse.json();
    }).then((simplePulls) => {
      const ids = simplePulls.map((pull) => pull.Id);
      var url = _pullDetailsWithIds(ids);

      return fetch(url).then((detailResponse) => {
        return detailResponse.json();
      }).then((Pulls) => {
        dispatch(receivePulls(Pulls));
      });
    });
  }
};

export const pageForwards = () => {
  return (dispatch, getState) => {
    var { page, numPulls } = getState();
    var page = page + numPulls;
    dispatch({
      type: actionTypes.CHANGE_PAGE,
      payload: page
    });
  }
};

export const pageBackwards = () => {
  return (dispatch, getState) => {
    var { page, numPulls } = getState();
    var page = page - numPulls;
    dispatch({
      type: actionTypes.CHANGE_PAGE,
      payload: page
    });
  };
};
