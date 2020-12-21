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
 * setCookie
 * @param {*} name
 * @param {*} value
 * @param {*} expire hours
 */
export const setCookie = (name, value, expire = 1) => {
  expire = expire * 60 * 60 * 1000;
  let exp = new Date();
  exp.setTime(exp.getTime() + expire);
  console.log('exp时间：', exp);
  document.cookie = `${name}=${escape(value)}; path=/;expires=${exp.toGMTString()}`;
}

/**
 * getCookie
 * @param {*} name
 */
export const getCookie = (name) => {
  let reg = new RegExp(`(^| )${name}=([^;]*)(;|$)`);
  let arr = document.cookie.match(reg);
  return arr ? unescape(arr[2]) : null;
}

/**
 * clearCookie
 * @param {*} name
 */
export const clearCookie = (name) => {
  setCookie(name, '', -1);
}

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
    return date.format('yyyyMMdd');
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
 * FindDeleteIndex  l2是l1删除某个元素后的新数组，找出l1被删除的元素的位置，没有返回-1
 * @params {*} l1
 * @params {*} l2
 * @returns index
 */
export const FindDeleteIndex = (l1, l2) => {
  let i = 0;
  while (i < l1.length && i < l2.length) {
    if (l1[i] !== l2[i]) {
      return i;
    }
    i++;
  }
  return i < l1.length ? i : -1;
}


/**
 * NoContainEle 找出了l2中l1不存在的元素
 * @params {*} l1
 * @params {*} l2
 * @returns list
 */
export const NoContainEle = (l1, l2) => {
  let obj = {};
  let l = [];
  for (let item of l1) {
    obj[item] = true;
  }
  for (let item of l2) {
    if (!obj[item]) {
      l.push(item);
    }
  }
  return l;
}

/**
 * DeleteElements 将l2的元素从l1中删除
 * @params {*} l1
 * @params {*} l2
 * @returns list
 */
export const DeleteElements = (l1, l2) => {
  return l1.filter(item => l2.indexOf(item) === -1);
}