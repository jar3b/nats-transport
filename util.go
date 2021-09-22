package nats_transport

// copy the values into protocol buffer
// struct
func copyMap(values map[string][]string) map[string]*Values {
	headerMap := make(map[string]*Values, 0)
	for k, v := range values {
		headerMap[k] = &Values{
			Arr: v,
		}
	}
	return headerMap
}
