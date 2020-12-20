/**
 * @file menu配置文件
 * @author 陈铭涛
 */

const getManagerMenu = context => {
    const menu = [
        {
            title: '预约管理',
            name: 'appointmentmanager',
            icon: 'md-time',
            to: {
                name: 'AppointmentManager'
            }
        },
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
            title: '校区',
            name: "campusmanager",
            icon: 'ios-school',
            to: {
                name: 'CampusManager',
            }
        },
        {
            title: '建筑',
            name: 'buildingmanager',
            icon: 'ios-pin',
            to: {
                name: 'BuildingManager'
            },
        },
        {
            title: '会议室',
            name: 'meetingmanager',
            icon: 'md-text',
            to: {
                name: 'MeetingManager'
            },
        },
    ];

    return menu;
};

export default getManagerMenu;
