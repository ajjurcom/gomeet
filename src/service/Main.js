import {WebHttp as http} from './index.js';

export default class Main {
    static login(data) {
        return http.post('/session', data);
    }
    static addUser(data) {
        return http.post('/user', data);
    }
    /*
     * 校区API
     */
    static addCampus(data) {
        return http.post('/campus', data);
    }
    static deleteCampus(id) {
        return http.delete(`/campus/${id}`);
    }
    static putCampus(data) {
        return http.put('/campus', data);
    }
    static getCampusByPage(onePageNum, page) {
        return http.get(`/campus/${onePageNum}/${page}`);
    }
    static getAllCampus() {
        return http.get('/campus');
    }
    /*
     * 建筑API
     */
    static addBuilding(data) {
        return http.post('/building', data);
    }
    static deleteBuilding(id) {
        return http.delete(`/building/${id}`);
    }
    static putBuilding(data) {
        return http.put('/building', data);
    }
    static getBuildingByCampusPage(onePageNum, page, campus_id) {
        return http.get(`/buildings/${onePageNum}/${page}?campus_id=${campus_id}`);
    }
    static getBuildingByID(id) {
        return http.get(`/building/${id}`);
    }
}