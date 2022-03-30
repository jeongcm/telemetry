package prometheus

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"strings"
	"telemetry"
	"time"
)

// Prometheus 는 prometheus에서 사용할 정보를 담는 구조체이다.
type Prometheus struct {
	client api.Client
}

// New 는 prometheus에 query를 요청하기 위한 client를 만드는 함수이다.
func New(name string) (telemetry.Telemetry, error) {
	services, err := registry.GetService(name)
	if err != nil && err != registry.ErrNotFound {
		logger.Errorf("Could not get service. cause: %v", err)
		return nil, err

	} else if err == registry.ErrNotFound || len(services) == 0 {
		logger.Errorf("Not found Service (%s)", name)
		return nil, telemetry.ErrNotFoundService
	}

	for _, svc := range services {
		for _, node := range svc.Nodes {
			var address string
			if !strings.HasPrefix(node.Address, "http://") {
				address = fmt.Sprint("http://", node.Address)
			}

			if client, err := api.NewClient(api.Config{Address: address}); err == nil {
				return &Prometheus{client: client}, nil
			}

			logger.Debugf("Connection Failed. %v", err)
		}
	}

	return nil, telemetry.ErrConnectionFailed
}

// promAPIQuery 는 prometheus에 query를 요청하는 함수이다.
func (p *Prometheus) promAPIQuery(query string) (model.Value, error) {
	v1api := v1.NewAPI(p.client)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, warnings, err := v1api.Query(ctx, query, time.Now().In(time.UTC))
	if err != nil {
		return nil, err
	}

	if len(warnings) > 0 {
		logger.Warnf("PromAPIQuery get warnings. cause: %v", warnings)
	}

	return result, nil
}

// NodeMeta 는 prometheus에서 node_meta 정보를 가져오는 함수이다.
func (p *Prometheus) NodeMeta(idList ...string) (map[string]telemetry.NodeMeta, error) {
	var query string
	if len(idList) > 0 {
		query = fmt.Sprintf("node_meta{container_label_com_docker_swarm_node_id=~\"%s\"}",
			strings.Join(idList, "|"))
	} else {
		query = "node_meta{}"
	}

	result, err := p.promAPIQuery(query)
	if err != nil {
		return nil, err
	}
	vectors := result.(model.Vector)
	nodes := make(map[string]telemetry.NodeMeta, len(vectors))
	for _, v := range vectors {
		metric := v.Metric
		key := fmt.Sprintf("%v", metric["instance"])
		nodes[key] = telemetry.NodeMeta{
			ID:        fmt.Sprintf("%v", metric["node_id"]),
			HostName:  fmt.Sprintf("%v", metric["node_name"]),
			IPAddress: fmt.Sprintf("%v", metric["node_ip"]),
		}
	}

	return nodes, nil
}

// NodeCPUCoreCnt 는 prometheus에서 node의 CPU Core 갯수를 가져오는 함수이다.
func (p *Prometheus) NodeCPUCoreCnt(key string) (int32, error) {
	query := fmt.Sprintf("count(count(node_cpu_seconds_total{instance=\"%s\"}) by (cpu))", key)
	result, err := p.promAPIQuery(query)
	if err != nil {
		return 0, err
	}
	vectors := result.(model.Vector)
	value := vectors[0].Value

	return int32(value), nil
}

// NodeCPUUsedRate 는 prometheus에서 node의 CPU 사용률을 가져오는 함수이다.
func (p *Prometheus) NodeCPUUsedRate(key string) (float32, error) {
	query := fmt.Sprintf("100 - (avg(irate(node_cpu_seconds_total{instance=~\"%s\",mode=\"idle\"}[5m]))"+
		" * 100)", key)
	result, err := p.promAPIQuery(query)
	if err != nil {
		return 0, err
	}
	vectors := result.(model.Vector)
	value := vectors[0].Value

	return float32(value), nil
}

// NodeMemTotalBytes 는 prometheus에서 node의 메모리 용량을 가져오는 함수이다.
func (p *Prometheus) NodeMemTotalBytes(key string) (int64, error) {
	query := fmt.Sprintf("sum(node_memory_MemTotal_bytes{instance=\"%s\"})", key)
	result, err := p.promAPIQuery(query)
	if err != nil {
		return 0, err
	}
	vectors := result.(model.Vector)
	value := vectors[0].Value

	return int64(value), nil
}

// NodeMemUsedBytes 는 prometheus에서 node의 메모리 사용량을 가져오는 함수이다.
func (p *Prometheus) NodeMemUsedBytes(key string) (int64, error) {
	query := fmt.Sprintf("sum(node_memory_MemTotal_bytes{instance=\"%s\"}) -"+
		"sum(node_memory_MemAvailable_bytes{instance=\"%s\"})", key, key)
	result, err := p.promAPIQuery(query)
	if err != nil {
		return 0, err
	}
	vectors := result.(model.Vector)
	value := vectors[0].Value

	return int64(value), nil
}

// NodeNetworkReceiveBytes 는 prometheus에서 node의 Network Input 사용량을 가져오는 함수이다.
func (p *Prometheus) NodeNetworkReceiveBytes(key string) (int64, error) {
	query := fmt.Sprintf("sum(rate(node_network_receive_bytes_total{instance=\"%s\"}[1m]))", key)
	result, err := p.promAPIQuery(query)
	if err != nil {
		return 0, err
	}
	vectors := result.(model.Vector)
	value := vectors[0].Value

	return int64(value), nil
}

// NodeNetworkTransmitBytes 는 prometheus에서 node의 Network Output 사용량을 가져오는 함수이다.
func (p *Prometheus) NodeNetworkTransmitBytes(key string) (int64, error) {
	query := fmt.Sprintf("sum(rate(node_network_transmit_bytes_total{instance=\"%s\"}[1m]))", key)
	result, err := p.promAPIQuery(query)
	if err != nil {
		return 0, err
	}
	vectors := result.(model.Vector)
	value := vectors[0].Value

	return int64(value), nil
}

// NodeFilesystemSizeBytes 는 prometheus에서 node의 스토리지 용량을 가져오는 함수이다.
func (p *Prometheus) NodeFilesystemSizeBytes(key string) (int64, error) {
	query := fmt.Sprintf("node_filesystem_size_bytes{instance=\"%s\"}", key)
	result, err := p.promAPIQuery(query)
	if err != nil {
		return 0, err
	}
	vectors := result.(model.Vector)
	value := vectors[0].Value

	return int64(value), nil
}

// NodeFilesystemUsedBytes 는 prometheus에서 node의 스토리지 사용량을 가져오는 함수이다.
func (p *Prometheus) NodeFilesystemUsedBytes(key string) (int64, error) {
	query := fmt.Sprintf("node_filesystem_size_bytes{instance=\"%s\"} -"+
		"node_filesystem_avail_bytes{instance=\"%s\"}", key, key)
	result, err := p.promAPIQuery(query)
	if err != nil {
		return 0, err
	}
	vectors := result.(model.Vector)
	value := vectors[0].Value

	return int64(value), nil
}

// ServiceMeta 는 prometheus에서 요청한 서비스가 실행중인 container의 기본정보를 가져오는 함수이다.
func (p *Prometheus) ServiceMeta(name string) ([]telemetry.ServiceMeta, error) {
	query := fmt.Sprintf("container_start_time_seconds{image!=\"\", "+
		"container_env_cdm_service_name=\"%s\"}", name)

	result, err := p.promAPIQuery(query)
	if err != nil {
		return nil, err
	}
	vectors := result.(model.Vector)
	services := make([]telemetry.ServiceMeta, 0, len(vectors))
	for _, v := range vectors {
		metric := v.Metric
		service := telemetry.ServiceMeta{
			ID:   fmt.Sprintf("%v", metric["container_label_com_docker_swarm_node_id"]),
			Name: fmt.Sprintf("%v", metric["container_env_cdm_service_name"]),
		}
		services = append(services, service)
	}

	return services, nil
}

// ServiceMemUsedBytes 는 prometheus에서 특정 서비스의 메모리 사용량을 가져오는 함수이다.
func (p *Prometheus) ServiceMemUsedBytes(params ...string) (int64, error) {
	var query string
	if len(params) > 1 {
		query = fmt.Sprintf("sum(container_memory_usage_bytes{image!=\"\", "+
			"container_env_cdm_service_name=\"%s\", container_label_com_docker_swarm_node_id=\"%s\"})",
			params[0], params[1])
	} else {
		query = fmt.Sprintf("sum(container_memory_usage_bytes{image!=\"\", "+
			"container_env_cdm_service_name=\"%s\"})", params[0])
	}

	result, err := p.promAPIQuery(query)
	if err != nil {
		return 0, err
	}
	vectors := result.(model.Vector)
	value := vectors[0].Value

	return int64(value), nil
}

// ServiceNetworkReceiveBytes 는 prometheus에서 특정 서비스의 Network Input 사용량을 가져오는 함수이다.
func (p *Prometheus) ServiceNetworkReceiveBytes(params ...string) (int64, error) {
	var query string
	if len(params) > 1 {
		query = fmt.Sprintf("sum(container_network_receive_bytes_total{image!=\"\", "+
			"container_env_cdm_service_name=\"%s\", container_label_com_docker_swarm_node_id=\"%s\"})",
			params[0], params[1])
	} else {
		query = fmt.Sprintf("sum(container_network_receive_bytes_total{image!=\"\", "+
			"container_env_cdm_service_name=\"%s\"})", params[0])
	}

	result, err := p.promAPIQuery(query)
	if err != nil {
		return 0, err
	}
	vectors := result.(model.Vector)
	value := vectors[0].Value

	return int64(value), nil
}

// ServiceNetworkTransmitBytes 는 prometheus에서 특정 서비스의 Network Output 사용량을 가져오는 함수이다.
func (p *Prometheus) ServiceNetworkTransmitBytes(params ...string) (int64, error) {
	var query string
	if len(params) > 1 {
		query = fmt.Sprintf("sum(container_network_transmit_bytes_total{image!=\"\", "+
			"container_env_cdm_service_name=\"%s\", container_label_com_docker_swarm_node_id=\"%s\"})",
			params[0], params[1])
	} else {
		query = fmt.Sprintf("sum(container_network_transmit_bytes_total{image!=\"\", "+
			"container_env_cdm_service_name=\"%s\"})", params[0])
	}

	result, err := p.promAPIQuery(query)
	if err != nil {
		return 0, err
	}
	vectors := result.(model.Vector)
	value := vectors[0].Value

	return int64(value), nil
}
