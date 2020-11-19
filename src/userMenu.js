/**
 * @file menu配置文件
 * @author 陈铭涛
 */


const getUserMenu = context => {
    const menu = [
        {
            title: '会议室预定',
            name: 'reserve-meeting',
            to: {
                name: 'ReserveMeeting'
            }
        },
        {
            title: '我的会议',
            name: 'my-meeting',
            to: {
                name: 'ReserveMeeting'
            }
        },
    ];

    return menu;
};

export default getUserMenu;