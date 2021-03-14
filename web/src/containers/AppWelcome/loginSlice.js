import { createSlice } from '@reduxjs/toolkit';

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
    console.log(login, password);
    setTimeout(() => {
        dispatch(setProgress(false));
        dispatch(setError('InvalidCredentials'));
    }, 1000);
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
