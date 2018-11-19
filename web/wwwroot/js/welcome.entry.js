import './utils/polyfills';

import Vue from 'vue';
import SignInFormComponent from './components/signin-form.vue';

var SignInForm = Vue.extend(SignInFormComponent);
new SignInForm().$mount('.login-form');