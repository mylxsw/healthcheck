import Vue from 'vue';
import Router from 'vue-router';
import Alerts from './views/Alerts';
import Healthchecks from './views/Healthchecks';

Vue.use(Router);

const routerPush = Router.prototype.push;
Router.prototype.push = function push(location) {
    return routerPush.call(this, location).catch(error => error)
}

export default new Router({
    routes: [
        {path: '/', component: Alerts},
        {path: '/healthchecks', component: Healthchecks},
    ]
});
