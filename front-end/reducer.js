import * as actionTypes from './actionTypes';

const initialState = {
  isFetching: false,
  page: 1,
  limit: 10
};

const reducer = (state = initialState, action) => {
  const payload = action.payload;
  
    switch(action.type) {
      case(actionTypes.REQUEST_PULLS):
        return {
          ...state,
          ...payload
        };
      case(actionTypes.RECEIVE_PULLS):
        return {
          ...state,
          ...payload
        };
      case(actionTypes.CHANGE_PAGE):
        return {
          ...state,
          page: payload
        };
      default:
        return state;
    }
}

export default reducer;
