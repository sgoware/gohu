package config

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4/storage"
)

// yaml 解析器
type emptyParser struct {
}

func (d *emptyParser) Parse(configContent interface{}) (map[string]interface{}, error) {

	return nil, nil
}

// CustomChangeListener 自定义客户端配置监控器
type CustomChangeListener struct {
}

func (c *CustomChangeListener) OnChange(changeEvent *storage.ChangeEvent) {
	fmt.Printf("Onchange(ChangeEvent)\nnamespace: %v,\n notificationId: %v\n", changeEvent.Namespace, changeEvent.NotificationID)
	for k, v := range changeEvent.Changes {
		fmt.Printf("key: %v\nChangeType: %v,\nOldValue: %v,\nNewValue: %v\n", k, v.ChangeType, v.OldValue, v.NewValue)
	}
}

func (c *CustomChangeListener) OnNewestChange(event *storage.FullChangeEvent) {
	fmt.Printf("OnNewestChange(FullChangeEvent)\nnamespace: %v,\n notificationId: %v\n", event.Namespace, event.NotificationID)
	for k, v := range event.Changes {
		fmt.Printf("key: %v, \nvalue: %v\n", k, v)
	}
}
