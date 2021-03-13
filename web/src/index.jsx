import React from 'react';
import ReactDOM from 'react-dom';
import App from './containers/AppMain/App';
import store from './containers/AppMain/store';
import { Provider } from 'react-redux';

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store}>
      <App />
    </Provider>
  </React.StrictMode>,
  document.getElementById('app')
);
