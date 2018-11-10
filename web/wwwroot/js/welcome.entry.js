require('./utils/polyfills');

var Vue = require('vue').default;
var VueRouter = require('vue-router').default;

Vue.use(VueRouter);

var WelcomeLayoutComponent = require('./components/welcome-layout/welcome-layout.component');
var WelcomeComponent = require('./components/welcome/welcome.component');
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

var app = new Vue({
    el: '#app',
    router: router,
    template: '<router-view></router-view>'
});