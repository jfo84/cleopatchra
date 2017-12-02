import React from 'react';
import { render } from 'react-dom';

import injectTapEventPlugin from 'react-tap-event-plugin';
import createSagaMiddleware from 'redux-saga';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import { cyan500 } from 'material-ui/styles/colors';

import { createStore, combineReducers, applyMiddleware } from 'redux';
import { Provider } from 'react-redux';

import createHistory from 'history/createBrowserHistory'
import { Route } from 'react-router';

import { ConnectedRouter, routerReducer, routerMiddleware, push } from 'react-router-redux';

import reducer from './reducer';

import Pulls from './components/Pulls';

const history = createHistory();

const middleware = routerMiddleware(history);
const sagaMiddleware = createSagaMiddleware();

const muiTheme = getMuiTheme({
  palette: {
    textColor: cyan500
  },
  appBar: {
    height: 50
  }
});

const store = createStore(
  combineReducers({
    ...reducer,
    ...applyMiddleware(sagaMiddleware),
    router: routerReducer
  }),
  applyMiddleware(middleware)
);

injectTapEventPlugin();

render(
  <MuiThemeProvider muiTheme={muiTheme}>
    <Provider store={store}>
      <ConnectedRouter history={history}>
        <div>
          <Route path="/repos/:repoId/pulls" component={Pulls}/>
        </div>
      </ConnectedRouter>
    </Provider>
  </MuiThemeProvider>,
  document.getElementById('root')
);
