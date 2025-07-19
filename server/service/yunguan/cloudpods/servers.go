package cloudpods

import (
	"ygpt/server/global"
	tt "yunion.io/x/onecloud/pkg/apis/monitor"
	"yunion.io/x/onecloud/pkg/mcclient/modules/monitor"
)

type ServerService struct {
}

// 磁盘写速率
func (i *ServerService) MonitorDiskWrite(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "agent_diskio",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "vm_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
					{
						Type: "tag",
						Params: []string{
							"vm_id",
						},
					},
					{
						Type: "tag",
						Params: []string{
							"name",
						},
					},
				},
				Selects: []tt.MetricQuerySelect{
					[]tt.MetricQueryPart{
						{
							Type: "field",
							Params: []string{
								"write_bps",
							},
						},
						{
							Type:   "mean",
							Params: []string{},
						},
						{
							Type: "alias",
							Params: []string{
								"磁盘写速率",
							},
						},
					},
				},
			},
		},
	}
	metricQueryInput := tt.MetricQueryInput{
		From:            time,
		Interval:        interval,
		MetricQuery:     metricQuery,
		Scope:           "system",
		Signature:       "",
		SkipCheckSeries: true,
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}
	data, err := monitor.UnifiedMonitorManager.PerformQuery(s, &metricQueryInput)
	if err != nil {
		return nil, err
	}
	return data.Interface(), err
}

// 磁盘读速率
func (i *ServerService) MonitorDiskRead(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "agent_diskio",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "vm_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
					{
						Type: "tag",
						Params: []string{
							"vm_id",
						},
					},
					{
						Type: "tag",
						Params: []string{
							"name",
						},
					},
				},
				Selects: []tt.MetricQuerySelect{
					[]tt.MetricQueryPart{
						{
							Type: "field",
							Params: []string{
								"read_bps",
							},
						},
						{
							Type:   "mean",
							Params: []string{},
						},
						{
							Type: "alias",
							Params: []string{
								"磁盘读速率",
							},
						},
					},
				},
			},
		},
	}
	metricQueryInput := tt.MetricQueryInput{
		From:            time,
		Interval:        interval,
		MetricQuery:     metricQuery,
		Scope:           "system",
		Signature:       "",
		SkipCheckSeries: true,
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}
	data, err := monitor.UnifiedMonitorManager.PerformQuery(s, &metricQueryInput)
	if err != nil {
		return nil, err
	}
	return data.Interface(), err
}

// 网络收包数
func (i *ServerService) MonitorPpsRecvData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "agent_net",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "vm_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
					{
						Type: "tag",
						Params: []string{
							"vm_id",
						},
					},
					{
						Type: "tag",
						Params: []string{
							"interface",
						},
					},
				},
				Selects: []tt.MetricQuerySelect{
					[]tt.MetricQueryPart{
						{
							Type: "field",
							Params: []string{
								"pps_recv",
							},
						},
						{
							Type:   "mean",
							Params: []string{},
						},
						{
							Type: "alias",
							Params: []string{
								"网络收包数/秒",
							},
						},
					},
				},
			},
		},
	}
	metricQueryInput := tt.MetricQueryInput{
		From:            time,
		Interval:        interval,
		MetricQuery:     metricQuery,
		Scope:           "system",
		Signature:       "",
		SkipCheckSeries: true,
		Unit:            false,
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}
	data, err := monitor.UnifiedMonitorManager.PerformQuery(s, &metricQueryInput)
	if err != nil {
		return nil, err
	}
	return data.Interface(), err
}

// 网络发包数
func (i *ServerService) MonitorPpsSentData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "agent_net",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "vm_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
					{
						Type: "tag",
						Params: []string{
							"vm_id",
						},
					},
					{
						Type: "tag",
						Params: []string{
							"interface",
						},
					},
				},
				Selects: []tt.MetricQuerySelect{
					[]tt.MetricQueryPart{
						{
							Type: "field",
							Params: []string{
								"pps_sent",
							},
						},
						{
							Type:   "mean",
							Params: []string{},
						},
						{
							Type: "alias",
							Params: []string{
								"网络发包数/秒",
							},
						},
					},
				},
			},
		},
	}
	metricQueryInput := tt.MetricQueryInput{
		From:            time,
		Interval:        interval,
		MetricQuery:     metricQuery,
		Scope:           "system",
		Signature:       "",
		SkipCheckSeries: true,
		Unit:            false,
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}
	data, err := monitor.UnifiedMonitorManager.PerformQuery(s, &metricQueryInput)
	if err != nil {
		return nil, err
	}
	return data.Interface(), err
}

// 网络入流量
func (i *ServerService) MonitorBpsRecvData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "agent_net",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "vm_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
					{
						Type: "tag",
						Params: []string{
							"vm_id",
						},
					},
					{
						Type: "tag",
						Params: []string{
							"interface",
						},
					},
				},
				Selects: []tt.MetricQuerySelect{
					[]tt.MetricQueryPart{
						{
							Type: "field",
							Params: []string{
								"bps_recv",
							},
						},
						{
							Type:   "mean",
							Params: []string{},
						},
						{
							Type: "alias",
							Params: []string{
								"网络入带宽",
							},
						},
					},
				},
			},
		},
	}
	metricQueryInput := tt.MetricQueryInput{
		From:            time,
		Interval:        interval,
		MetricQuery:     metricQuery,
		Scope:           "system",
		Signature:       "",
		SkipCheckSeries: true,
		Unit:            false,
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}
	data, err := monitor.UnifiedMonitorManager.PerformQuery(s, &metricQueryInput)
	if err != nil {
		return nil, err
	}
	return data.Interface(), err
}

// 网络出流量
func (i *ServerService) MonitorBpsSentData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "agent_net",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "vm_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
					{
						Type: "tag",
						Params: []string{
							"vm_id",
						},
					},
					{
						Type: "tag",
						Params: []string{
							"interface",
						},
					},
				},
				Selects: []tt.MetricQuerySelect{
					[]tt.MetricQueryPart{
						{
							Type: "field",
							Params: []string{
								"bps_sent",
							},
						},
						{
							Type:   "mean",
							Params: []string{},
						},
						{
							Type: "alias",
							Params: []string{
								"网络出带宽",
							},
						},
					},
				},
			},
		},
	}
	metricQueryInput := tt.MetricQueryInput{
		From:            time,
		Interval:        interval,
		MetricQuery:     metricQuery,
		Scope:           "system",
		Signature:       "",
		SkipCheckSeries: true,
		Unit:            false,
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}
	data, err := monitor.UnifiedMonitorManager.PerformQuery(s, &metricQueryInput)
	if err != nil {
		return nil, err
	}
	return data.Interface(), err
}

// cpu使用率
func (i *ServerService) MonitorCpuData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "agent_cpu",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "vm_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
					{
						Type: "tag",
						Params: []string{
							"vm_id",
						},
					},
				},
				Selects: []tt.MetricQuerySelect{
					[]tt.MetricQueryPart{
						{
							Type: "field",
							Params: []string{
								"usage_active",
							},
						},
						{
							Type:   "mean",
							Params: []string{},
						},
						{
							Type: "alias",
							Params: []string{
								"CPU使用率",
							},
						},
					},
				},
			},
		},
	}
	metricQueryInput := tt.MetricQueryInput{
		From:            time,
		Interval:        interval,
		MetricQuery:     metricQuery,
		Scope:           "system",
		Signature:       "",
		SkipCheckSeries: true,
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}
	data, err := monitor.UnifiedMonitorManager.PerformQuery(s, &metricQueryInput)
	if err != nil {
		return nil, err
	}
	return data.Interface(), err
}

// 内存使用率
func (i *ServerService) MonitorMemData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "agent_mem",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "vm_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
					{
						Type: "tag",
						Params: []string{
							"vm_id",
						},
					},
				},
				Selects: []tt.MetricQuerySelect{
					[]tt.MetricQueryPart{
						{
							Type: "field",
							Params: []string{
								"used_percent",
							},
						},
						{
							Type:   "mean",
							Params: []string{},
						},
						{
							Type: "alias",
							Params: []string{
								"内存使用率",
							},
						},
					},
				},
			},
		},
	}

	metricQueryInput := tt.MetricQueryInput{
		From:            time,
		Interval:        interval,
		MetricQuery:     metricQuery,
		Scope:           "system",
		Signature:       "",
		SkipCheckSeries: true,
		Unit:            false,
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}
	data, err := monitor.UnifiedMonitorManager.PerformQuery(s, &metricQueryInput)
	if err != nil {
		return nil, err
	}
	return data.Interface(), err
}

// 磁盘使用率
func (i *ServerService) MonitorDiskData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "agent_disk",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "vm_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
					{
						Type: "tag",
						Params: []string{
							"vm_id",
						},
					},
					{
						Type: "tag",
						Params: []string{
							"device",
						},
					},
					{
						Type: "tag",
						Params: []string{
							"path",
						},
					},
					{
						Type: "tag",
						Params: []string{
							"fstype",
						},
					},
				},
				Selects: []tt.MetricQuerySelect{
					[]tt.MetricQueryPart{
						{
							Type: "field",
							Params: []string{
								"used_percent",
							},
						},
						{
							Type:   "mean",
							Params: []string{},
						},
						{
							Type: "alias",
							Params: []string{
								"磁盘使用率",
							},
						},
					},
				},
			},
		},
	}

	metricQueryInput := tt.MetricQueryInput{
		From:            time,
		Interval:        interval,
		MetricQuery:     metricQuery,
		Scope:           "system",
		Signature:       "",
		SkipCheckSeries: true,
		Unit:            false,
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}
	data, err := monitor.UnifiedMonitorManager.PerformQuery(s, &metricQueryInput)
	if err != nil {
		return nil, err
	}
	return data.Interface(), err
}

// GPU使用率
func (i *ServerService) MonitorGpuData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "agent_nvidia_smi",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "vm_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
					{
						Type: "tag",
						Params: []string{
							"index",
						},
					},
				},
				Selects: []tt.MetricQuerySelect{
					[]tt.MetricQueryPart{
						{
							Type: "field",
							Params: []string{
								"utilization_gpu",
							},
						},
						{
							Type:   "mean",
							Params: []string{},
						},
						{
							Type: "alias",
							Params: []string{
								"GPU使用率",
							},
						},
					},
				},
			},
		},
	}

	metricQueryInput := tt.MetricQueryInput{
		From:            time,
		Interval:        interval,
		MetricQuery:     metricQuery,
		Scope:           "system",
		Signature:       "",
		SkipCheckSeries: true,
		Unit:            false,
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}
	data, err := monitor.UnifiedMonitorManager.PerformQuery(s, &metricQueryInput)
	if err != nil {
		return nil, err
	}
	return data.Interface(), err
}

// GPU显存使用率
func (i *ServerService) MonitorGpuMemData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "agent_nvidia_smi",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "vm_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
					{
						Type: "tag",
						Params: []string{
							"index",
						},
					},
				},
				Selects: []tt.MetricQuerySelect{
					[]tt.MetricQueryPart{
						{
							Type: "field",
							Params: []string{
								"utilization_memory",
							},
						},
						{
							Type:   "mean",
							Params: []string{},
						},
						{
							Type: "alias",
							Params: []string{
								"GPU显存使用率",
							},
						},
					},
				},
			},
		},
	}

	metricQueryInput := tt.MetricQueryInput{
		From:            time,
		Interval:        interval,
		MetricQuery:     metricQuery,
		Scope:           "system",
		Signature:       "",
		SkipCheckSeries: true,
		Unit:            false,
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}
	data, err := monitor.UnifiedMonitorManager.PerformQuery(s, &metricQueryInput)
	if err != nil {
		return nil, err
	}
	return data.Interface(), err
}

// GPU温度
func (i *ServerService) MonitorGpuTemperatureData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "agent_nvidia_smi",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "vm_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
					{
						Type: "tag",
						Params: []string{
							"index",
						},
					},
				},
				Selects: []tt.MetricQuerySelect{
					[]tt.MetricQueryPart{
						{
							Type: "field",
							Params: []string{
								"temperature_gpu",
							},
						},
						{
							Type:   "mean",
							Params: []string{},
						},
						{
							Type: "alias",
							Params: []string{
								"GPU温度",
							},
						},
					},
				},
			},
		},
	}

	metricQueryInput := tt.MetricQueryInput{
		From:            time,
		Interval:        interval,
		MetricQuery:     metricQuery,
		Scope:           "system",
		Signature:       "",
		SkipCheckSeries: true,
		Unit:            false,
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}
	data, err := monitor.UnifiedMonitorManager.PerformQuery(s, &metricQueryInput)
	if err != nil {
		return nil, err
	}
	return data.Interface(), err
}
