package storage

type StorageApiGroup struct {
	ContainerEvsApi
	ContainerNasApi
	VmEvsApi
	VmNasApi
	StoSnapshotApi
}
