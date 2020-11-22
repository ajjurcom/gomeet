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
    static putUserState(id, state) {
        return http.put(`/user_state/${id}?state=${state}`);
    }
    static putUserPwd(data) {
        return http.put('/user_password', data);
    }
    static getUser(id) {
        return http.get(`/user/${id}`);
    }
    static getUserOptions(role) {
        return http.get(`/user_options?role=${role}`);
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
    static searchUsers(obj) {
        return http.get(`/users?searchWay=${obj.searchWay}&keyword=${obj.keyword}`)
    }
    static getGroupMembers(id) {
        return http.get(`/members/${id}`);
    }
    /*
     * 用户组API
     */
    static addGrouop(data) {
        return http.post('/user_group', data);
    }
    static deleteGroup(id) {
        return http.delete(`/user_group/${id}`);
    }
    static putGroupName(data) {
        return http.put('/user_name', data);
    }
    static putGroupMember(data) {
        return http.put('/user_member', data);
    }
    static getGroupsByPage(obj) {
        return http.get(`/user_group/${obj.onePageNum}/${obj.page}?creator=${obj.creator}`);
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
    /*
     * 预约API
     */
    static getScheduleOptions() {
        return http.get('/schedule_options');
    }
}