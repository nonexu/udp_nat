package misc

const (
	CLUSTER_TYPE = 1
	CLUSTER_MSG  = 2
)

type AddressInfo struct {
	Id      int32
	Address string
}

type ClusterInfo struct {
	Id   int32
	Node []*AddressInfo
}

type Msg struct {
	Type int16
	Data string
}
