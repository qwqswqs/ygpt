package cloudpods

import (
	"ygpt/server/global"
	tt "yunion.io/x/onecloud/pkg/apis/monitor"
	"yunion.io/x/onecloud/pkg/mcclient/modules/monitor"
)

type BareHostService struct {
}

// 磁盘写速率
func (i *BareHostService) MonitorDiskWrite(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "diskio",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "host_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
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
								"write_bps",
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
func (i *BareHostService) MonitorDiskRead(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "diskio",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "host_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
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
								"read_bps",
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

// 网络入流量
func (i *BareHostService) MonitorBpsRecvData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "net",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "host_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
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
								"bps_recv",
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
func (i *BareHostService) MonitorBpsSentData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "net",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "host_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
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
								"bps_sent",
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
func (i *BareHostService) MonitorCpuData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "cpu",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "host_id",
						Operator: "=",
						Value:    id,
					},
					{
						Key:      "cpu",
						Operator: "=",
						Value:    "cpu-total",
					},
				},
				GroupBy: []tt.MetricQueryPart{
					{
						Type: "tag",
						Params: []string{
							"cpu",
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
func (i *BareHostService) MonitorMemData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "mem",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "host_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{},
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
					}, []tt.MetricQueryPart{
						{
							Type: "field",
							Params: []string{
								"free_percent",
							},
						},
						{
							Type:   "mean",
							Params: []string{},
						},
						{
							Type: "alias",
							Params: []string{
								"内存空闲率",
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
func (i *BareHostService) MonitorDiskData(id, time, interval string) (interface{}, error) {
	metricQuery := []*tt.AlertQuery{
		{
			Model: tt.MetricQuery{
				Measurement: "disk",
				Tags: []tt.MetricQueryTag{
					{
						Key:      "host_id",
						Operator: "=",
						Value:    id,
					},
				},
				GroupBy: []tt.MetricQueryPart{
					{
						Type: "tag",
						Params: []string{
							"path",
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
								"used_percent",
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
