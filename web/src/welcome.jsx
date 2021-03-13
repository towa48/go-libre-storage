import React from 'react';
import ReactDOM from 'react-dom';
import App from './containers/AppWelcome/AppWelcome';
import store from './containers/AppWelcome/store';
import { Provider } from 'react-redux';

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store}>
      <App />
    </Provider>
  </React.StrictMode>,
  document.getElementById('app')
);
