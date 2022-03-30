package prometheus

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	p, err := New(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, p)

	//not found service
	p2, err2 := New("unknown service")
	assert.Error(t, err2)
	assert.Nil(t, p2)
}

func TestPrometheus_NodeMeta(t *testing.T) {
	p, err := New(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := p.NodeMeta()
	assert.NoError(t, err)
	assert.NotNil(t, nodes)

	node, err2 := p.NodeMeta("xmxzh74mzu19nhbl4iwdnc79i")
	assert.NoError(t, err2)
	assert.NotNil(t, node)
}

func TestPrometheus_NodeCPUCoreCnt(t *testing.T) {
	p, err := New(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := p.NodeMeta()
	if err != nil {
		t.Fatal(err)
	}
	for key := range nodes {
		core, err := p.NodeCPUCoreCnt(key)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, core)
	}

}

func TestPrometheus_NodeCPUUsedRate(t *testing.T) {
	p, err := New(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := p.NodeMeta()
	if err != nil {
		t.Fatal(err)
	}
	for key := range nodes {
		rate, err := p.NodeCPUUsedRate(key)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, rate)
	}

}

func TestPrometheus_NodeMemTotalBytes(t *testing.T) {
	p, err := New(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := p.NodeMeta()
	if err != nil {
		t.Fatal(err)
	}
	for key := range nodes {
		bytes, err := p.NodeMemTotalBytes(key)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, bytes)
	}
}

func TestPrometheus_NodeMemUsedBytes(t *testing.T) {
	p, err := New(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := p.NodeMeta()
	if err != nil {
		t.Fatal(err)
	}
	for key := range nodes {
		bytes, err := p.NodeMemUsedBytes(key)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, bytes)
	}
}

func TestPrometheus_NodeNetworkReceiveBytes(t *testing.T) {
	p, err := New(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := p.NodeMeta()
	if err != nil {
		t.Fatal(err)
	}
	for key := range nodes {
		bytes, err := p.NodeNetworkReceiveBytes(key)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, bytes)
	}
}

func TestPrometheus_NodeNetworkTransmitBytes(t *testing.T) {
	p, err := New(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := p.NodeMeta()
	if err != nil {
		t.Fatal(err)
	}
	for key := range nodes {
		bytes, err := p.NodeNetworkTransmitBytes(key)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, bytes)
	}
}

func TestPrometheus_NodeFilesystemSizeBytes(t *testing.T) {
	p, err := New(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := p.NodeMeta()
	if err != nil {
		t.Fatal(err)
	}
	for key := range nodes {
		bytes, err := p.NodeFilesystemSizeBytes(key)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, bytes)
	}
}

func TestPrometheus_NodeFilesystemUsedBytes(t *testing.T) {
	p, err := New(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := p.NodeMeta()
	if err != nil {
		t.Fatal(err)
	}
	for key := range nodes {
		bytes, err := p.NodeFilesystemUsedBytes(key)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, bytes)
	}
}

func TestPrometheus_ServiceMeta(t *testing.T) {
	p, err := New(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}

	_, err = p.ServiceMeta(constant.BackingServiceTelemetry)
	assert.NoError(t, err)
}

func TestPrometheus_ServiceMemUsedBytes(t *testing.T) {
	p, err := New(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}

	services, err := p.ServiceMeta(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}

	for _, svc := range services {
		TotalMemUsedBytes, err := p.ServiceMemUsedBytes(svc.Name)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, TotalMemUsedBytes)

		MemUsedBytes, err := p.ServiceMemUsedBytes(svc.Name, svc.ID)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, MemUsedBytes)
	}
}

func TestPrometheus_ServiceNetworkReceiveBytes(t *testing.T) {
	p, err := New(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}

	services, err := p.ServiceMeta(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}

	for _, svc := range services {
		TotalNetworkReceiveBytes, err := p.ServiceNetworkReceiveBytes(svc.Name)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, TotalNetworkReceiveBytes)

		NetworkReceiveBytes, err := p.ServiceNetworkReceiveBytes(svc.Name, svc.ID)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, NetworkReceiveBytes)
	}
}

func TestPrometheus_ServiceNetworkTransmitBytes(t *testing.T) {
	p, err := New(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}

	services, err := p.ServiceMeta(constant.BackingServiceTelemetry)
	if err != nil {
		t.Fatal(err)
	}

	for _, svc := range services {
		TotalNetworkTransmitBytes, err := p.ServiceNetworkTransmitBytes(svc.Name)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, TotalNetworkTransmitBytes)

		NetworkTransmitBytes, err := p.ServiceNetworkTransmitBytes(svc.Name, svc.ID)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, NetworkTransmitBytes)
	}
}
