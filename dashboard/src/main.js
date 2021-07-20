import '@babel/polyfill'
import 'mutationobserver-shim'
import Vue from 'vue'
import './plugins/axios'
import './plugins/bootstrap-vue'
import App from './App.vue'
import router from './router'

import { BootstrapVueIcons } from 'bootstrap-vue'

import DateTime from "./components/DateTime";
import HumanTime from "./components/HumanTime";

Vue.use(BootstrapVueIcons);

Vue.component('DateTime', DateTime);
Vue.component('HumanTime', HumanTime);

Vue.config.productionTip = false;

const errorHandler = (error) => {
    console.log(error);
}

Vue.config.errorHandler = errorHandler;
Vue.prototype.$throw = (error) => errorHandler(error, this);

Vue.prototype.QueryArgs = (route, name) => {
    return route.query[name] !== undefined ? route.query[name] : '';
}

/**
 * @return {string}
 */
Vue.prototype.ParseError = function (error) {
    if (error.response !== undefined) {
        if (error.response.status === 401) {
            this.$throw("access-denied");
        } 

        if (error.response.data !== undefined) {
            return error.response.data.error;
        }
    }

    return error.toString();
};

Vue.prototype.ToastSuccess = function (message) {
    this.$bvToast.toast(message, {
        title: 'OK',
        variant: 'success',
        autoHideDelay: 3000,
        toaster: 'b-toaster-top-center',
    });
};

Vue.prototype.ToastError = function (message) {
    this.$bvToast.toast(this.ParseError(message), {
        title: 'Oops',
        variant: 'danger',
        autoHideDelay: 3000,
        toaster: 'b-toaster-top-center',
    });
};

Vue.prototype.SuccessBox = function (message, cb) {
    cb = cb || function () {};
    this.$bvModal.msgBoxOk(message, {
        title: 'Success',
        centered: true,
        okVariant: 'success',
        headerClass: 'p-2 border-bottom-0',
        footerClass: 'p-2 border-top-0',
    }).then(cb);
};

Vue.prototype.ErrorBox = function (message, cb) {
    cb = cb || function () {};
    this.$bvModal.msgBoxOk(this.ParseError(message), {
        centered: true,
        title:'Oops',
        okVariant: 'danger',
        headerClass: 'p-2 border-bottom-0',
        footerClass: 'p-2 border-top-0',
    }).then(cb);
};

new Vue({
    router,
    render: h => h(App)
}).$mount('#app');
