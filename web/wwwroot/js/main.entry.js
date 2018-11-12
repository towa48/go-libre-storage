import './utils/polyfills';

import Vue from 'vue';
import VueRouter from 'vue-router';

Vue.use(VueRouter);

import IndexLayoutComponent from './components/index-layout.vue';
import IndexComponent from './components/index.vue';
var NotFoundComponent = { template: '<h2>Not found</h2>' }

var routes = [{
    path: '/index',
    component: IndexLayoutComponent,
    children: [{
        path: '',
        name: 'index',
        component: IndexComponent
    }]
}, {
    path: '/',
    redirect: '/index'
}, {
    path: '*',
    component: IndexLayoutComponent,
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
    template: '<router-view></router-view>'
});