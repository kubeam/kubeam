package models

/*ClusterList struct describes responses with description of clusters*/
type ClusterList struct {
	Description string
	Clusters    map[string]map[string]string
}
