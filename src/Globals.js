import MainAPI from '@/service/Main';


import viewDesign from 'view-design';
import 'view-design/dist/styles/iview.css';


export default {
    install(Vue, option) {
        Vue.use(viewDesign);
        Vue.prototype.$service = {
            MainAPI
        };
        Vue.prototype.$changeLoadingState = function (key, value) {
            this.loading[key] = value;
        };
    }
};