import { createSlice } from '@reduxjs/toolkit';
import { HttpStatus, $http } from '../../utils/http';
import { $location } from '../../utils/location';
import routes from './routes';

const initialState = {
    login: '',
    password: '',
    inProgress: false,
    isPasswordShown: false,
    error: null,
}

export const loginFormSlice = createSlice({
    name: 'loginForm',
    initialState: initialState,
    reducers: {
        setLogin: (state, action) => {
            state.login = action.payload;
        },
        setPassword: (state, action) => {
            state.password = action.payload;
        },
        setProgress: (state, action) => {
            state.inProgress = action.payload;
        },
        togglePasswordVisibility: (state) => {
            state.isPasswordShown = !state.isPasswordShown;
        },
        setError: (state, action) => {
            if (state.error === action.payload) {
                return;
            }
            state.error = action.payload;
        },
        clearError: (state) => {
            state.error = initialState.error;
        }
    },
});

/*
 * Actions
 */
export const {
    setLogin,
    setPassword,
    setProgress,
    togglePasswordVisibility,
    setError,
    clearError,
} = loginFormSlice.actions;

/*
 * Effects (thunks)
 */
export const doLoginAsync = (login, password) => dispatch => {
    dispatch(clearError());
    dispatch(setProgress(true));

    const data = {
        username: login,
        password: password,
        rememberMe: false
    };

    $http.post(routes.signIn, data).then(function(resp) {
        switch(resp.status) {
            case HttpStatus.OK:
                resp.json().then(function(data) {
                    $location.url(data.Url || '/');
                });
                break;
            case HttpStatus.BAD_REQUEST:
                resp.json().then(function(data) {
                    dispatch(setProgress(false));
                    dispatch(setError('InvalidCredentials'));
                });
                break;
            case HttpStatus.SERVER_ERROR:
                dispatch(setProgress(false));
                dispatch(setError('ServerError'));
                break;
        }
    }).catch(function(err) {
        dispatch(setProgress(false));
        dispatch(setError('UnknownError'));
    });
}

/*
 * Selectors
 */
export const selectLogin = state => state.loginForm.login;
export const selectPassword = state => state.loginForm.password;
export const selectProgress = state => state.loginForm.inProgress;
export const selectPasswordVisibility = state => state.loginForm.isPasswordShown;
export const selectError = state => state.loginForm.error;

/*
 * Reducer for composition
 */
export default loginFormSlice.reducer;
