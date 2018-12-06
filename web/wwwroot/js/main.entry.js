import './utils/polyfills';

import Vue from 'vue';
import VueRouter from 'vue-router';
import Fragment from 'vue-fragment';

Vue.use(VueRouter);
Vue.use(Fragment.Plugin);

import IndexLayoutComponent from './components/index-layout.vue';
import DocumentsComponent from './components/documents.vue';
import SharedActivityComponent from './components/shared-activity.vue';
var NotFoundComponent = { template: '<h2>Not found</h2>' }

var routes = [{
    path: '/documents',
    component: IndexLayoutComponent,
    children: [{
        path: '',
        name: 'documents',
        component: DocumentsComponent
    }]
}, {
    path: '/shared/activity',
    component: IndexLayoutComponent,
    children: [{
        path: '',
        name: 'sharedActivity',
        component: SharedActivityComponent
    }]
}, {
    path: '/',
    redirect: '/documents'
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