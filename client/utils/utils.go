package utils

import (
	"chapterRoom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

// 这里将方法关联到结构体中
type Transfer struct {
	// 分析应该有哪些字段
	Conn net.Conn   // 连接
	Buf  [8096]byte //传输时使用的缓冲
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	//buffer := make([]byte, 8096)
	//fmt.Println("读取客户端发送的数据...")
	// conn.Read在conn没有关闭的情况下，才会阻塞
	// 如果客户端关闭了conn则就不会阻塞
	// 读长度
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		fmt.Println("in readPkg conn.Read err=", err)
		return
	}

	//fmt.Println("读到的buffer(ok)=", this.Buf[:4])
	//根据buf[:4],转为一个uint32类型
	// 为了知道读多少字节
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])
	// 根据pkgLen读取信息内容
	// 读 写进buffer中
	n1, err1 := this.Conn.Read(this.Buf[:pkgLen])
	//fmt.Println("buffer[:pkgLen](ok) = ", this.Buf[:pkgLen])
	// n1 = pkgLen才是正常的
	if n1 != int(pkgLen) || err1 != nil {
		fmt.Println("conn.Read fail err=", err)
		//err = errors.New("read pkg body error")
		return
	}

	// 把pkgLen反序列化Message类型
	// 技术就是一层窗户纸 &mes
	// 注意是传递mes引用,并且这是函数返回值，并没有传递进来
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		//err = errors.New("json.Unmarshal err")
		return
	}
	//fmt.Println("in readPkg mes = ",mes)
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	// 先发送一个长度给对方，此时server和client已经不那么重要了
	var pkgLen uint32
	pkgLen = uint32(len(data))
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	// 发送长度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}

	// 发送data
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	return
}
