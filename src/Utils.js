import { Message } from 'view-design';

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
    var dd = new Date();
    dd.setDate(dd.getDate() + afterDay);
    return dd;
}

/**
 * ShowDateFormat
 * @params {*} date
 * @returns "10/27 周二" 
 */
export const ShowDateFormat = (date) => {
    const month = (date.getMonth() + 1) < 10 ? "0" + (date.getMonth() + 1) : (date.getMonth() + 1); //获取当前月份的日期，不足10补0
    const d = date.getDate() < 10 ? "0" + date.getDate() : date.getDate(); //获取当前几号，不足10补0
    const day = "周" + "日一二三四五六".charAt(date.getDay());
    return month + "/" + d + " " + day;
}

/**
 * GetNumFromScale
 * @params {*} str
 * @returns num 
 */
export const GetNumFromScale = (str) => {
    return str.replace(/[^0-9]/ig, "") || 0;
}