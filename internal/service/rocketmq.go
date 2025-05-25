package service

import (
	"fmt"

	ons20190214 "github.com/alibabacloud-go/ons-20190214/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

// RocketMQService handles RocketMQ operations
type RocketMQService struct {
	client *ons20190214.Client
}

// NewRocketMQService creates a new RocketMQService
func NewRocketMQService(client *ons20190214.Client) *RocketMQService {
	return &RocketMQService{client: client}
}

// RocketMQInstance represents a RocketMQ instance
type RocketMQInstance struct {
	InstanceId     string `json:"instanceId"`
	InstanceName   string `json:"instanceName"`
	InstanceType   int32  `json:"instanceType"`
	InstanceStatus int32  `json:"instanceStatus"`
	RegionId       string `json:"regionId"`
	CreateTime     int64  `json:"createTime"`
	ReleaseTime    int64  `json:"releaseTime"`
	Remark         string `json:"remark"`
	ServiceVersion int32  `json:"serviceVersion"`
}

// RocketMQTopic represents a RocketMQ topic
type RocketMQTopic struct {
	Topic       string `json:"topic"`
	MessageType int32  `json:"messageType"`
	InstanceId  string `json:"instanceId"`
	CreateTime  int64  `json:"createTime"`
	UpdateTime  int64  `json:"updateTime"`
	Remark      string `json:"remark"`
	Status      int32  `json:"status"`
	Perm        int32  `json:"perm"`
}

// RocketMQGroup represents a RocketMQ consumer group
type RocketMQGroup struct {
	GroupId    string `json:"groupId"`
	GroupType  string `json:"groupType"`
	InstanceId string `json:"instanceId"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
	Remark     string `json:"remark"`
}

// FetchInstances retrieves all RocketMQ instances
func (s *RocketMQService) FetchInstances() ([]RocketMQInstance, error) {
	request := &ons20190214.OnsInstanceInServiceListRequest{}

	response, err := s.client.OnsInstanceInServiceList(request)
	if err != nil {
		return nil, fmt.Errorf("fetching RocketMQ instances: %w", err)
	}

	var instances []RocketMQInstance
	if response.Body != nil && response.Body.Data != nil && response.Body.Data.InstanceVO != nil {
		for _, inst := range response.Body.Data.InstanceVO {
			instance := RocketMQInstance{
				InstanceId:     tea.StringValue(inst.InstanceId),
				InstanceName:   tea.StringValue(inst.InstanceName),
				InstanceType:   tea.Int32Value(inst.InstanceType),
				InstanceStatus: tea.Int32Value(inst.InstanceStatus),
				RegionId:       "", // Field not available in API response
				CreateTime:     tea.Int64Value(inst.CreateTime),
				ReleaseTime:    tea.Int64Value(inst.ReleaseTime),
				Remark:         "", // Field not available in API response
				ServiceVersion: 0,  // Field not available in API response
			}
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

// FetchTopics retrieves all topics for a specific RocketMQ instance
func (s *RocketMQService) FetchTopics(instanceId string) ([]RocketMQTopic, error) {
	request := &ons20190214.OnsTopicListRequest{
		InstanceId: tea.String(instanceId),
	}

	response, err := s.client.OnsTopicList(request)
	if err != nil {
		return nil, fmt.Errorf("fetching topics for instance %s: %w", instanceId, err)
	}

	var topics []RocketMQTopic
	if response.Body != nil && response.Body.Data != nil && response.Body.Data.PublishInfoDo != nil {
		for _, topic := range response.Body.Data.PublishInfoDo {
			topicInfo := RocketMQTopic{
				Topic:       tea.StringValue(topic.Topic),
				MessageType: tea.Int32Value(topic.MessageType),
				InstanceId:  tea.StringValue(topic.InstanceId),
				CreateTime:  tea.Int64Value(topic.CreateTime),
				UpdateTime:  0, // Field not available in API response
				Remark:      tea.StringValue(topic.Remark),
				Status:      0, // Field not available in API response
				Perm:        0, // Field not available in API response
			}
			topics = append(topics, topicInfo)
		}
	}

	return topics, nil
}

// FetchGroups retrieves all consumer groups for a specific RocketMQ instance
func (s *RocketMQService) FetchGroups(instanceId string) ([]RocketMQGroup, error) {
	request := &ons20190214.OnsGroupListRequest{
		InstanceId: tea.String(instanceId),
	}

	response, err := s.client.OnsGroupList(request)
	if err != nil {
		return nil, fmt.Errorf("fetching groups for instance %s: %w", instanceId, err)
	}

	var groups []RocketMQGroup
	if response.Body != nil && response.Body.Data != nil && response.Body.Data.SubscribeInfoDo != nil {
		for _, group := range response.Body.Data.SubscribeInfoDo {
			groupInfo := RocketMQGroup{
				GroupId:    tea.StringValue(group.GroupId),
				GroupType:  tea.StringValue(group.GroupType),
				InstanceId: tea.StringValue(group.InstanceId),
				CreateTime: tea.Int64Value(group.CreateTime),
				UpdateTime: tea.Int64Value(group.UpdateTime),
				Remark:     tea.StringValue(group.Remark),
			}
			groups = append(groups, groupInfo)
		}
	}

	return groups, nil
}
