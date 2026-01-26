// internal/runtime/metadata.go
package runtime

import (
	"net"
	"os"

	"github.com/rs/zerolog/log"
)

// Metadata contém informações sobre o ambiente de execução
type Metadata struct {
	Hostname    string
	PodName     string
	PodIP       string
	NodeName    string
	Namespace   string
	AppVersion  string
	Environment string
}

// NewMetadata inicializa os metadados do runtime (pod/container)
func NewMetadata(appVersion, environment string) *Metadata {
	hostname, _ := os.Hostname()

	// Em Kubernetes, o hostname geralmente é o pod name
	podName := hostname
	if envPodName := os.Getenv("POD_NAME"); envPodName != "" {
		podName = envPodName
	}

	metadata := &Metadata{
		Hostname:    hostname,
		PodName:     podName,
		PodIP:       getPodIP(),
		NodeName:    os.Getenv("NODE_NAME"),
		Namespace:   os.Getenv("POD_NAMESPACE"),
		AppVersion:  appVersion,
		Environment: environment,
	}

	log.Info().
		Str("hostname", metadata.Hostname).
		Str("pod_name", metadata.PodName).
		Str("pod_ip", metadata.PodIP).
		Str("node_name", metadata.NodeName).
		Str("namespace", metadata.Namespace).
		Str("app_version", metadata.AppVersion).
		Str("environment", metadata.Environment).
		Msg("Runtime metadata initialized")

	return metadata
}

// getPodIP obtém o IP do pod/container
func getPodIP() string {
	// Tenta primeiro a variável de ambiente (Kubernetes Downward API)
	if podIP := os.Getenv("POD_IP"); podIP != "" {
		return podIP
	}

	// Fallback: detecta o IP da interface de rede principal
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return "unknown"
}

// GetIdentifier retorna um identificador único do pod (para logs)
func (m *Metadata) GetIdentifier() string {
	if m.Namespace != "" {
		return m.Namespace + "/" + m.PodName
	}
	return m.PodName
}

// ToLogFields retorna os campos para logging estruturado
func (m *Metadata) ToLogFields() map[string]interface{} {
	fields := make(map[string]interface{})

	if m.Hostname != "" {
		fields["hostname"] = m.Hostname
	}
	if m.PodName != "" {
		fields["pod_name"] = m.PodName
	}
	if m.PodIP != "" {
		fields["pod_ip"] = m.PodIP
	}
	if m.NodeName != "" {
		fields["node_name"] = m.NodeName
	}
	if m.Namespace != "" {
		fields["namespace"] = m.Namespace
	}
	if m.AppVersion != "" {
		fields["app_version"] = m.AppVersion
	}
	if m.Environment != "" {
		fields["environment"] = m.Environment
	}

	return fields
}
