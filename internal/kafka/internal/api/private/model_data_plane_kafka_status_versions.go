/*
 * Kafka Service Fleet Manager
 *
 * Kafka Service Fleet Manager APIs that are used by internal services e.g kas-fleetshard operators.
 *
 * API version: 1.6.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package private

// DataPlaneKafkaStatusVersions Version information related to a Kafka cluster
type DataPlaneKafkaStatusVersions struct {
	Kafka    string `json:"kafka,omitempty"`
	Strimzi  string `json:"strimzi,omitempty"`
	KafkaIbp string `json:"kafkaIbp,omitempty"`
}
