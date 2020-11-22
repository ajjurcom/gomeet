import {Message} from 'view-design';

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
    for (let i = 0; i < length-1; i++) {
        str = str + intList[i] + ",";
    }
    str += intList[length-1];
    return str;
}
