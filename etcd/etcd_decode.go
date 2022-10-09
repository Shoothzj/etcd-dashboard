package etcd

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"k8s.io/kubectl/pkg/scheme"
)

func decode(component, namespace string, content []byte) (string, error) {
	if component == "Kubernetes" {
		return decodeKubernetes(content)
	}
	return string(content), nil
}

func decodeKubernetes(content []byte) (string, error) {
	decoder := scheme.Codecs.UniversalDeserializer()
	obj, _, err := decoder.Decode(content, nil, nil)
	if err != nil {
		logrus.Errorf("Failed to decode kubernetes object: %v", err)
		return "", err
	}

	jsonData, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		logrus.Errorf("Failed to marshal kubernetes object: %v", err)
		return "", err
	}
	return string(jsonData), nil
}
