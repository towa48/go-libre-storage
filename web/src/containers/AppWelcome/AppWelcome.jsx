import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import './AppWelcome.scss';

import LoginForm from '../../components/LoginForm/LoginForm';
import {
    setLogin,
    setPassword,
    doLoginAsync,
    togglePasswordVisibility,
    selectLogin,
    selectPassword,
    selectProgress,
    selectPasswordVisibility,
    selectError,
} from './loginSlice';

function App() {
    const login = useSelector(selectLogin);
    const password = useSelector(selectPassword);
    const isPasswordShown = useSelector(selectPasswordVisibility);
    const progress = useSelector(selectProgress);
    const error = useSelector(selectError);
    const dispatch = useDispatch();

    return (
        <LoginForm
          login={login}
          password={password}
          inProgress={progress}
          isPasswordShown={isPasswordShown}
          error={error}
          onChangeLogin={(login) => dispatch(setLogin(login))}
          onChangePassword={(password) => dispatch(setPassword(password))}
          onSubmit={(login, password) => dispatch(doLoginAsync(login, password))}
          onTogglePasswordVisibility={() => dispatch(togglePasswordVisibility())}/>
    );
}

export default App;
