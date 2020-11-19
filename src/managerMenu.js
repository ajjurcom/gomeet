/**
 * @file menu配置文件
 * @author 陈铭涛
 */

const getManagerMenu = context => {
    const menu = [
        {
            title: '用户管理',
            name: "usermanager",
            icon: 'logo-octocat',
            to: {
                name: 'UserManager',
                query: {
                    state: 'verify_user'
                }
            }
        },
        {
            title: '用户组',
            name: 'group',
            icon: 'md-chatboxes',
            children: [
                {
                    title: '管理用户组',
                    name: 'group-manager'
                },
                {
                    title: '新增用户组',
                    name: 'group-add'
                },
            ]
        },
        {
            title: '校区',
            name: "campus",
            icon: 'ios-school',
            children: [
                {
                    title: '管理校区',
                    name: 'campusmanager',
                    to: {
                        name: 'CampusManager'
                    }
                },
                {
                    title: '新增校区',
                    name: 'campusadd',
                    to: {
                        name: 'CampusAdd'
                    }
                }
            ]
        },
        {
            title: '建筑',
            name: 'building',
            icon: 'ios-pin',
            children: [
                {
                    title: '管理建筑',
                    name: 'buildingmanager',
                    to: {
                        name: 'BuildingManager'
                    }
                },
                {
                    title: '新增建筑',
                    name: 'buildingadd',
                    to: {
                        name: 'BuildingAdd'
                    }
                }
            ]
        },
        {
            title: '会议室',
            name: 'meeting',
            icon: 'md-text',
            children: [
                {
                    title: '管理会议室',
                    name: 'meetingmanager',
                    to: {
                        name: 'MeetingManager'
                    }
                },
                {
                    title: '新增会议室',
                    name: 'meetingadd',
                    to: {
                        name: 'MeetingAdd'
                    }
                },
            ]
        },
    ];

    return menu;
};

export default getManagerMenu;
