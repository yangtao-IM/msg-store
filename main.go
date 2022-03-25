package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

//记录每个用户入库的最新一条msgid,[uk][msgid]
var CursorCache = make(map[int]int)
var cacheMutex sync.Mutex

type MsgDTO struct {
	W    *http.ResponseWriter
	done chan struct{}
	Msg  *Msg
}

var Queue = make(chan *MsgDTO)

type Msg struct {
	Uk      int    `json:"uk"`
	ToUk    int    `json:"to_uk"`
	Id      int64  `json:"id"`
	Content string `json:"content"`
}

type SaveMsgRsp struct {
	MsgId   int64 `json:"msgid"`
	ErrCode int64 `json:"err_code"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/msg/save", func(w http.ResponseWriter, r *http.Request) {
		SaveMsgHandler(w, r)
	})
	go consumMsg()
	http.ListenAndServe(":8888", mux)
}

func SaveMsgHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	msg := &Msg{}
	err = json.Unmarshal(b, msg)
	if err != nil {
		fmt.Println(err.Error())
	}

	done := make(chan struct{})
	fmt.Printf("req:%+v", string(b))
	//后续可以做sync.Pool的优化
	dto := &MsgDTO{
		W:    &w, //入库成功以后返回msgid
		done: done,
		Msg:  msg,
	}
	Queue <- dto
	<-done //等待入库成功，写完msgid以后再返回，http的请求在退出handler方法以后会返回
}

func consumMsg() {
	for {
		for dto := range Queue {
			msg := dto.Msg
			//生成消息id
			id := time.Now().UnixNano()
			msg.Id = id
			rsp := &SaveMsgRsp{
				MsgId:   id,
				ErrCode: 200,
			}
			b, err := json.Marshal(rsp)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			(*dto.W).Write(b)

			//这里需要加锁，如果能保证一个消费协程一个map，则不用加锁
			cacheMutex.Lock()
			CursorCache[msg.Uk] = int(msg.Id)
			cacheMutex.Unlock()

			dto.done <- struct{}{}
		}
	}
}
