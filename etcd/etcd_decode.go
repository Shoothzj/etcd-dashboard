package etcd

func decode(component, namespace string, content []byte) (string, error) {
	return string(content), nil
}
