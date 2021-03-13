import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import './AppWelcome.css';

import LoginForm from '../../components/LoginForm/LoginForm';
import {
    setLogin,
    setPassword,
    doLoginAsync,
    selectLogin,
    selectPassword,
    selectProgress
} from './loginSlice';

function App() {
    const login = useSelector(selectLogin);
    const password = useSelector(selectPassword);
    const progress = useSelector(selectProgress);
    const dispatch = useDispatch();

    return (
        <LoginForm
          login={login}
          password={password}
          inProgress={progress}
          onChangeLogin={(login) => dispatch(setLogin(login))}
          onChangePassword={(password) => dispatch(setPassword(password))}
          onSubmit={(login, password) => dispatch(doLoginAsync(login, password))}/>
    );
}

export default App;
