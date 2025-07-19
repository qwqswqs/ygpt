package cloudpods

type CloudpodsRouter struct {
	SystemImageRouter
	DockerImageRouter
	ContainerRouter
	VirtualMachineRouter
	BareMetalRouter
	JumpRouter
	CloudpodsResInfoRouter
}
