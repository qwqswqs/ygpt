package compute

type ComputeRouter struct {
	ServersRouter
	ScalingGroupRouter
	ScalingConfigRouter
	ScalingPolicyRouter
}
