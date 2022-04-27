package kubernetes

type RemovePodsData struct {
	Namespace string `json:"namespace"  binding:"required"`
	PodName   string `json:"podName"  binding:"required"`
}
