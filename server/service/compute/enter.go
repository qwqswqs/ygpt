package compute

import "ygpt/server/service/compute/udp"

type ComputeServiceGroup struct {
	//minio.MinioServiceGroup
	udp.UdpServiceGroup
	//websocket.WebsocketServiceGroup
}
