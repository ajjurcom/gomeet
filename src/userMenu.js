/**
 * @file menu配置文件
 * @author 陈铭涛
 */


const getUserMenu = context => {
    const menu = [
        {
            title: '会议室预定',
            name: 'reservemeeting',
            to: {
                name: 'ReserveMeeting'
            }
        },
        {
            title: '我的会议',
            name: 'mymeeting',
            to: {
                name: 'ReserveMeeting'
            }
        },
    ];

    return menu;
};

export default getUserMenu;