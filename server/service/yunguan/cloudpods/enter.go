package cloudpods

type CloudpodsServiceGroup struct {
	VirtualMachineService
	BareMetalService
	ContainerService
	ImageService
	ResInfoService
	ServerService
	BareHostService
}
