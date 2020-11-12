import {WebHttp as http} from './index.js';

export default class Main {
    /*
     * 用户API
     */
    static login(data) {
        return http.post('/session', data);
    }
    static addUser(data) {
        return http.post('/user', data);
    }
    static putUser(data) {
        return http.put('/user', data);
    }
    static putUserPwd(data) {
        return http.put('/user_password', data);
    }
    static getUser(id) {
        return http.get(`/user/${id}`);
    }
    static getAllStateList() {
        return http.get('/user_states');
    }
    static getUserInfo(id) {
        return http.get(`/user/${id}`);
    }
    static getUsersByPage(obj) {
        return http.get(`/users/${obj.onePageNum}/${obj.page}?state=${obj.state}`);
    }
    static deleteUser(id) {
        return http.delete(`/user/${id}`);
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
    /*
     * 会议室API
     */
    static addMeeting(data) {
        return http.post('/meeting', data);
    }
    static deleteMeeting(id) {
        return http.delete(`/meeting/${id}`);
    }
    static putMeeting(data) {
        return http.put('/meeting', data);
    }
    static getMeetingByID(id) {
        return http.get(`/meeting/${id}`);
    }
    static getMeetingOptions() {
        return http.get('/meeting_options');
    }
    static getAllBuildingsByCampus(campus_id) {
        return http.get(`/campus_buildings/${campus_id}`);
    }
    static getBuildingLayer(campus_id) {
        return http.get(`/campus_layer/${campus_id}`);
    }
    static getMeetingsByPage(onePageNum, page, building_id) {
        return http.get(`/meetings/${onePageNum}/${page}?building_id=${building_id}`);
    }
}