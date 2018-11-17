import './utils/polyfills';

import Vue from 'vue';
import VueRouter from 'vue-router';

Vue.use(VueRouter);

import WelcomeLayoutComponent from './components/welcome-layout.vue';
import WelcomeComponent from './components/welcome.vue';
var NotFoundComponent = { template: '<h2>Not found</h2>' }

var routes = [{
    path: '/welcome',
    component: WelcomeLayoutComponent,
    children: [{
        path: '',
        name: 'welcome',
        component: WelcomeComponent
    }]
}, {
    path: '/',
    redirect: '/welcome'
}, {
    path: '*',
    component: WelcomeLayoutComponent,
    children: [{
        path: '',
        name: 'notfound',
        component: NotFoundComponent
    }]
}];

var router = new VueRouter({
    routes: routes
})

new Vue({
    el: '#app',
    router: router,
    //template: '<router-view></router-view>'
});