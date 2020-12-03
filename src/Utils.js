import { Message } from 'view-design';

// 给Date加格式化方法
Date.prototype.format = function (fmt) {
    var o = {
      "M+": this.getMonth() + 1, //月份
      "d+": this.getDate(), //日
      "h+": this.getHours(), //小时
      "m+": this.getMinutes(), //分
      "s+": this.getSeconds(), //秒
      "q+": Math.floor((this.getMonth() + 3) / 3), //季度
      S: this.getMilliseconds(), //毫秒
    };
    if (/(y+)/.test(fmt)) {
      fmt = fmt.replace(
        RegExp.$1,
        (this.getFullYear() + "").substr(4 - RegExp.$1.length)
      );
    }
    for (var k in o) {
      if (new RegExp("(" + k + ")").test(fmt)) {
        fmt = fmt.replace(
          RegExp.$1,
          RegExp.$1.length == 1 ? o[k] : ("00" + o[k]).substr(("" + o[k]).length)
        );
      }
    }
    return fmt;
  };

/**
 * showMessage 全局显示提示消息
 * @params {*} type
 * @params {*} content
 */
export const showMessage = (type, content) => {
    return Message[type] && Message[type]({
        content,
        duration: 3,
        closable: true
    });
};

/**
 * setLocalStorage 设置localStorage
 * @params {*} name
 * @params {*} value
 */
export const setLocalStorage = (name, value) => {
    window.localStorage.setItem(name, value);
}

/**
 * getLocalStorage
 * @params {*} name
 */
export const getLocalStorage = (name) => {
    return window.localStorage.getItem(name);
}

/**
 * removeLocalStorage
 * @params {*} name
 */
export const removeLocalStorage = (name) => {
    return window.localStorage.removeItem(name);
}

/**
 * intArrayToStr
 * @params {*} intList
 */

export const intArrayToStr = (intList) => {
    let str = ""
    let length = intList.length
    if (length === 0) {
        return str;
    }
    for (let i = 0; i < length - 1; i++) {
        str = str + intList[i] + ",";
    }
    str += intList[length - 1];
    return str;
}

/**
 * GetDateObj
 * @params {*} afterDay
 */
export const GetDateObj = (afterDay) => {
    let dd = new Date();
    dd.setDate(dd.getDate() + afterDay);
    return dd;
}

/**
 * ReserveFormat 预定表展示的时间, 包含周几，不包含年
 * @params {*} date
 * @returns "10/27 周二" 
 */
export const ReserveFormat = (date) => {
    return date.format('MM/dd') + " " + "周" + "日一二三四五六".charAt(date.getDay());
}

/**
 * DateFormat 提交的日期格式, 包含年月日
 * @params {*} date
 * @returns "10/27/2020" 
 */
export const DateFormat = (date) => {
    return date.format('MM/dd/yyyy');
}

/**
 * GetNumFromScale
 * @params {*} str
 * @returns num 
 */
export const GetNumFromScale = (str) => {
    return str.replace(/[^0-9]/ig, "") || 0;
}


/**
 * GetNumberArr
 * @params {*} a    a < b
 * @params {*} b
 * @returns list
 */
export const GetNumberArr = (a, b) => {
    return Array.from(Array( b - a + 1 )).map(( e, i ) => a + i);
}


/**
 * GetNumFromScale
 * @function myTask任务函数 
 * @way {hour | minutes} str
 */
// export const OnTimeTask = (myTask, way) => {
//     const date = new Date();
//     const dateM = date.getMinutes();
//     const dateS = date.getSeconds();
//     const dateMS = 1000 - date.getMilliseconds();
//     let interval, after;
//     if (way && way === 'minutes') {
//         interval = 60000;
//         after = interval - dateS * 1000 - dateMS;
//     } else {
//         interval = 3600000;
//         after = interval - dateM * 60000 - dateS * 1000 - dateMS;
//     }
//     console.log('after 毫秒 -> ', after);
//     setTimeout(function () {
//         myTask();
//         setInterval(() => {
//             myTask();
//         }, interval);
//     }, after);
// }