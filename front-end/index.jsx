import React from 'react';
import injectTapEventPlugin from 'react-tap-event-plugin';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import createSagaMiddleware from 'redux-saga';

import { render } from 'react-dom';
import { createStore, applyMiddleware } from 'redux';
import { Provider } from 'react-redux';
import { cyan500 } from 'material-ui/styles/colors';

import App from './components/App';
import reducer from './reducer';

injectTapEventPlugin();

const sagaMiddleware = createSagaMiddleware();

const muiTheme = getMuiTheme({
  palette: {
    textColor: cyan500
  },
  appBar: {
    height: 50
  }
});

export const store = createStore(
  reducer,
  applyMiddleware(sagaMiddleware)
);

export default render(
  <MuiThemeProvider muiTheme={muiTheme}>
    <Provider store={store}>
      <App/>
    </Provider>
  </MuiThemeProvider>,
  document.getElementById('app')
);
