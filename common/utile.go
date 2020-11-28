package common

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
)

// Encryption 使用sha265算法加密密码
func Encryption(password string) string {
	h := sha256.New()
	io.WriteString(h, password)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// MemberListToStr 将成员列表转化为字符串保存
func MemberListToStr(idList []string) string {
	var tmp strings.Builder
	l := len(idList)
	if l == 0 {
		return ""
	}
	for i := 0; i < l-1; i++ {
		tmp.WriteString(idList[i])
		tmp.WriteString(",")
	}
	tmp.WriteString(idList[l-1])
	return tmp.String()
}

// MemberStrToList 将成员字符串转为列表
func MemberStrToList(idStr string) []string {
	str := strings.Trim(idStr, ",")
	if str == "" {
		return []string{}
	}
	return strings.Split(str, ",")
}

// RemoveDuplicateEle 删除字符串slice中的重复项
func RemoveDuplicateEle(l []string) []string {
	m := make(map[string]struct{})
	result := make([]string, 0, len(l))
	for _, v := range l {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

/* DiffMember 计算出idList、idList2的差异元素
 * 1. deleteList: idList中有，idList2中没有的
 * 2. newList: idList2中有，idList中没有的
 */
func DiffMember(oldList []string, newList []string) (deletedList []string, addList []string) {
	oldDict := map[string]bool{}
	newDict := map[string]bool{}
	for _, v := range oldList {
		oldDict[v] = true
	}
	for _, v := range newList {
		newDict[v] = true
	}

	for _, v := range oldList {
		if _, exists := newDict[v]; !exists {
			deletedList = append(deletedList, v)
		}
	}
	for _, v := range newList {
		if _, exists := oldDict[v]; !exists {
			addList = append(addList, v)
		}
	}
	return
}

// 某个值在切片中的索引, 不存在返回-1
func StrIndexOf(strList []string, val string) int {
	for i, v := range strList {
		if v == val {
			return i
		}
	}
	return -1
}