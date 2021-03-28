import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import StyledEngineProvider from '@material-ui/core/StyledEngineProvider';
import CssBaseline from '@material-ui/core/CssBaseline';
import { ThemeProvider } from '@material-ui/core/styles';
import App from './containers/AppMain/App';
import store from './containers/AppMain/store';
import theme from './theme';

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store}>
      <StyledEngineProvider injectFirst>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <App />
        </ThemeProvider>
      </StyledEngineProvider>
    </Provider>
  </React.StrictMode>,
  document.getElementById('app')
);
