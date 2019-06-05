package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "strconv"
    "reflect"
    "crypto/md5"
    "encoding/hex"
)

func main() {
    dir := "./"
    renLogFileName := ".renlog"
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        fmt.Printf("Failed to read dir %s, the error is %v\n", dir, err)
        return
    }
    renameLog := ""
    num := 1;
    for _, f := range files {
        fileName := f.Name()
        if fileName == renLogFileName {
            continue
        }
        // 代表取文件名第一个字符
        if fileName[:1] == "." {
            fmt.Printf("Skip hidden file %s\n", fileName)
            continue
        }
        var imageExtensions = []string{".jpg",".png"}
        // 只修改当前文件目录下文件
        if !f.IsDir() {
            dotIndex := strings.LastIndex(fileName, ".")
            //判断文件是一个合法的常规文件
            if dotIndex != -1 && dotIndex != 0 {
                extensionName := fileName[dotIndex:]
                ok , _ := in_array(extensionName,imageExtensions)
                if !ok{
                    continue
                }
                fmt.Printf("file: %s\n",fileName)
                // go里面没有do while循环体，只能用这种写法
                // int转string
                newFileName := strconv.Itoa(num)
                // 注意：++ 和 --只能作为语句使用，不能作为表达式
                num++
                newFileName += extensionName
                // 如果新文件名被占用，就把它先修改成随机的，这种逻辑还是最省事了
                fmt.Printf("%s 文件即将被修改成 %s", fileName, newFileName)
                if file_exists(newFileName) {
                    // 如果当前文件和新文件同名，直接跳过
                    if fileName == newFileName{
                        fmt.Printf("，因为同名，直接跳过 \n")
                        continue
                    }
                    fmt.Printf("，%s 文件名已经被占用\n", newFileName)
                    tmpFileName := GetMD5Hash(fileName)
                    tmpFileName += extensionName
                    os.Rename(newFileName, tmpFileName)
                    fmt.Printf("\n把 %s 临时修改成 %s\n", fileName, tmpFileName)
                }
                fmt.Printf("，转化 %s 为 %s\n", fileName, newFileName)
                err = os.Rename(fileName, newFileName)
                if err != nil {
                    fmt.Printf("Failed to rename file %s to %s, the error is %v\n", fileName, newFileName, err)
                    continue
                }
                renameLog += fmt.Sprintf("%s\t%s\n", fileName, newFileName)
            }
        }
    }
    fmt.Printf(renameLog)
    ioutil.WriteFile(renLogFileName, []byte(renameLog), 0666)
}

func in_array(val interface{}, array interface{}) (exists bool, index int) {
    exists = false
    index = -1

    switch reflect.TypeOf(array).Kind() {
    case reflect.Slice:
        s := reflect.ValueOf(array)

        for i := 0; i < s.Len(); i++ {
            if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
                index = i
                exists = true
                return
            }
        }
    }

    return
}
// 判断所给路径文件/文件夹是否存在
func file_exists(path string) bool {
    _, err := os.Stat(path)    //os.Stat获取文件信息
    if err != nil {
        if os.IsExist(err) {
            return true
        }
        return false
    }
    return true
}
func GetMD5Hash(text string) string {
    return GetByteMD5Hash([]byte(text))
}
func GetByteMD5Hash(content []byte) string {
    hasher := md5.New()
    hasher.Write(content)
    return hex.EncodeToString(hasher.Sum(nil))
}
