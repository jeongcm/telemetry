package telemetry

// NodeMeta 는 Docker Swarm Node 의 정보를 담는 구조체이다.
type NodeMeta struct {
	ID        string
	HostName  string
	IPAddress string
}

// ServiceMeta 는 Docker Service 의 정보를 담는 구조체이다.
type ServiceMeta struct {
	ID   string
	Name string
}

// Telemetry 는 서비스나 노드의 리소스 정보를 받아오는 함수를 정의한 인터페이스이다.
type Telemetry interface {
	NodeMeta(...string) (map[string]NodeMeta, error)
	NodeCPUCoreCnt(string) (int32, error)
	NodeCPUUsedRate(string) (float32, error)
	NodeMemTotalBytes(string) (int64, error)
	NodeMemUsedBytes(string) (int64, error)
	NodeNetworkReceiveBytes(string) (int64, error)
	NodeNetworkTransmitBytes(string) (int64, error)
	NodeFilesystemSizeBytes(string) (int64, error)
	NodeFilesystemUsedBytes(string) (int64, error)
	ServiceMeta(string) ([]ServiceMeta, error)
	ServiceMemUsedBytes(...string) (int64, error)
	ServiceNetworkReceiveBytes(...string) (int64, error)
	ServiceNetworkTransmitBytes(...string) (int64, error)
}
