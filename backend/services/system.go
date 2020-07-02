package services

import (
	"crawunit/constants"
	"crawunit/database"
	"crawunit/entity"
	"crawunit/model"
	"crawunit/utils"
	"encoding/json"
)

// 系统信息 chan 映射
var SystemInfoChanMap = utils.NewChanMap()

// 从远端获取系统信息
func GetRemoteSystemInfo(nodeId string) (sysInfo entity.SystemInfo, err error) {
	// 发送消息
	msg := entity.NodeMessage{
		Type:   constants.MsgTypeGetSystemInfo,
		NodeId: nodeId,
	}

	// 序列化
	msgBytes, _ := json.Marshal(&msg)
	if _, err := database.RedisClient.Publish("nodes:"+nodeId, utils.BytesToString(msgBytes)); err != nil {
		return entity.SystemInfo{}, err
	}

	// 通道
	ch := SystemInfoChanMap.ChanBlocked(nodeId)

	// 等待响应，阻塞
	sysInfoStr := <-ch

	// 反序列化
	if err := json.Unmarshal([]byte(sysInfoStr), &sysInfo); err != nil {
		return sysInfo, err
	}

	return sysInfo, nil
}

// 获取系统信息
func GetSystemInfo(nodeId string) (sysInfo entity.SystemInfo, err error) {
	if IsMasterNode(nodeId) {
		sysInfo, err = model.GetLocalSystemInfo()
	} else {
		sysInfo, err = GetRemoteSystemInfo(nodeId)
	}
	return
}



