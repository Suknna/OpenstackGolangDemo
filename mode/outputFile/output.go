package outputfile

import (
	"fmt"
	"os"
)

/*
1.生成云主机信息和节点信息到文件中
2.检测生成目录下是否包含同名文件如果包含则重命名旧文件
*/

// 将处理后的字符串输出到日志中
func OutputFile(data string, filePath string) error {
	// 写入
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open file %s failed, err: %s", filePath, err.Error())
	}
	defer file.Close()
	file.WriteString(data)
	return nil
}
