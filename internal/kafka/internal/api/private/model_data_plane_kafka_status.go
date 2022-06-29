/*
 * Kafka Service Fleet Manager
 *
 * Kafka Service Fleet Manager APIs that are used by internal services e.g kas-fleetshard operators.
 *
 * API version: 1.6.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package private

// DataPlaneKafkaStatus Schema of the status object for a Kafka cluster
type DataPlaneKafkaStatus struct {
	// The status conditions of a Kafka cluster
	Conditions []DataPlaneClusterUpdateStatusRequestConditions `json:"conditions,omitempty"`
	Capacity   DataPlaneKafkaStatusCapacity                    `json:"capacity,omitempty"`
	Versions   DataPlaneKafkaStatusVersions                    `json:"versions,omitempty"`
	// Routes created for a Kafka cluster
	Routes         *[]DataPlaneKafkaStatusRoutes `json:"routes,omitempty"`
	AdminServerURI string                        `json:"adminServerURI,omitempty"`
}
