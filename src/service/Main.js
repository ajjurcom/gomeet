import {WebHttp as http} from './index.js';

export default class Main {
    static login(data) {
        return http.post('/session', data);
    }
    static addUser(data) {
        return http.post('/user', data);
    }
}