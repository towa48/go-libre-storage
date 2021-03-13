import { createSlice } from '@reduxjs/toolkit';

export const loginFormSlice = createSlice({
  name: 'loginForm',
  initialState: {
    login: '',
    password: '',
    inProgress: false
  },
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
  },
});

/*
 * Actions
 */
export const { setLogin, setPassword, setProgress } = loginFormSlice.actions;

/*
 * Effects (thunks)
 */
export const doLoginAsync = (login, password) => dispatch => {
    dispatch(setProgress(true));
    console.log(login, password);
    setTimeout(() => {
        dispatch(setProgress(false));
    }, 1000);
}

/*
 * Selectors
 */
export const selectLogin = state => state.loginForm.login;
export const selectPassword = state => state.loginForm.password;
export const selectProgress = state => state.loginForm.inProgress;

/*
 * Reducer for composition
 */
export default loginFormSlice.reducer;
