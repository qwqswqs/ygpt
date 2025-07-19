package storage

type StorageRouterGroup struct {
	StoSnapshotRouter
	ContainerEvsRouter
	ContainerNasRouter
	VmEvsRouter
	VmNasRouter
}
