package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/kubeam/kubeam/common"
	"github.com/kubeam/kubeam/models"
	"gopkg.in/yaml.v2"
)

/*NewDBClient establishes a new redis database connection and returns
the client connection object*/
func NewDBClient() *redis.Client {
	redisHost, err := common.Config.GetString("redis/host", "localhost")
	redisPort, err := common.Config.GetInt("redis/port", 6379)
	redisPassword, err := common.Config.GetString("redis/password", "")
	common.LogInfo.Println(redisHost)
	common.LogInfo.Println(redisPort)
	common.LogInfo.Println(redisPassword)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", redisHost, redisPort),
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return (client)
}

/*DBClientReserveCluster updates the redis cache and allocates the cluster
to the incoming resource*/
func DBClientReserveCluster(client *redis.Client, app string, env string, key string, val []byte, t time.Duration) error {
	err := client.Set(fmt.Sprintf("%v-%v-%v", app, env, key), val, t).Err()
	//first slice then convert to string (string is a read-only slice of bytes)
	common.LogInfo.Printf("%v-%v-%v [%v]", app, env, key, string(val[:]))
	if err != nil {
		return err
	}
	common.LogInfo.Printf("Created reservation for %v-%v-%v [%v] with TTL of %v", app, env, key, string(val[:]), 300)
	return nil
}

/*DBClientFindAndReserve is a wrapper to find a free cluster and update the
cache of used clusters*/
func DBClientFindAndReserve(client *redis.Client, app string, env string, ttl time.Duration) (string, error) {

	var clusters models.ClusterList

	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("clusters/%v-%v-clusterlist.yaml", app, env))
	if err != nil {
		return "", fmt.Errorf("Could not find a cluster definition for application: %v env: %v", app, env)
	}
	err = yaml.Unmarshal(yamlFile, &clusters)
	check(err)

	common.LogInfo.Printf("Description: %#v\n", clusters.Description)
	for key, value := range clusters.Clusters {
		fmt.Println("Checking Cluster", key)
		fmt.Println("Checking value", value)
		val, err := client.Get(fmt.Sprintf("%v-%v-%v", app, env, key)).Result()
		if err == redis.Nil {
			common.LogInfo.Printf("Found available cluster [%v] for you.", key)
			decodedValue, _ := json.Marshal(value)
			ttl, err := time.ParseDuration("1500s")
			if err != nil {
				defaultTTL, _ := strconv.Atoi(value["default_ttl"])
				defaultTTLParsed := time.Duration(defaultTTL)
				ttl = time.Duration(defaultTTLParsed * time.Second)
			}
			err = DBClientReserveCluster(client, app, env, string(key), decodedValue, ttl)
			if err != nil {
				return "", err
			}
			return string(key), nil
		} else if err != nil {
			common.LogError.Printf("Failed to query redis for key [%v-%v-%v]", app, env, key)
		} else {
			common.LogInfo.Printf("Cluster %v has a reservation valid for 0 seconds %v", key, val)
		}
	}
	common.LogInfo.Printf("No clusters available for reservation for appliation %v environment %v", app, env)
	return "", fmt.Errorf("No clusters available for reservation for appliation %v environment %v", app, env)

}

/*DBGetClusterReservation ...*/
func DBGetClusterReservation(client *redis.Client, app string, env string, cluster string) (string, error) {

	var clusters models.ClusterList

	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("clusters/%v-%v-clusterlist.yaml", app, env))
	if err != nil {
		return "", fmt.Errorf("Could not find a cluster definition for application: %v env: %v", app, env)
	}
	err = yaml.Unmarshal(yamlFile, &clusters)
	check(err)

	common.LogInfo.Printf("Description: %#v\n", clusters.Description)
	//gotReservation := false
	var output bytes.Buffer
	output.WriteString("{")
	val, err := client.Get(fmt.Sprintf("%v-%v-%v", app, env, cluster)).Result()
	if err == redis.Nil {
		common.LogInfo.Printf("Cluster %v is free", cluster)
	} else if err != nil {
		common.LogError.Printf("Failed to query redis for key [%v-%v-%v]", app, env, cluster)
	} else {
		common.LogInfo.Printf("Cluster %v has a reservation valid for 0 seconds %v", cluster, val)

		out := map[string]interface{}{}
		json.Unmarshal([]byte(val), &out)

		out["application"] = app
		out["environment"] = env
		out["cluster"] = cluster
		keyExp := client.TTL(fmt.Sprintf("%v-%v-%v", app, env, cluster))
		out["ttl"] = keyExp.String()

		outputJSON, _ := json.Marshal(out)
		output.WriteString(string(outputJSON))
	}
	output.WriteString("}")
	common.LogInfo.Printf("List Of clusters [%v]", output)
	return output.String(), err

}

/*DBClientGetSingleCluster wrapper to get resource using a given cluster*/
func DBClientGetSingleCluster(client *redis.Client, app string, env string, cluster string) (string, error) {
	ret, err := DBClientListClusters(client, app, env, cluster, false)
	return ret, err
}

/*DBClientGetSingleClusterDetail wrapper to get details of a resource using a
given cluster*/
func DBClientGetSingleClusterDetail(client *redis.Client, app string, env string, cluster string) (string, error) {
	ret, err := DBClientListClusters(client, app, env, cluster, true)
	return ret, err
}

/*DBClientGetAllClusters is a wrapper to list the clusters in use and the
resources using the clusters*/
func DBClientGetAllClusters(client *redis.Client, app string, env string) (string, error) {
	ret, err := DBClientListClusters(client, app, env, "", false)
	return ret, err
}

/*DBClientGetAllClustersDetail is a wrapper to get details used clusters and
resources using the clusters*/
func DBClientGetAllClustersDetail(client *redis.Client, app string, env string) (string, error) {
	ret, err := DBClientListClusters(client, app, env, "", true)
	return ret, err
}

/*DBClientListClusters fetches the details of clusters in use and resources
using those clusters*/
func DBClientListClusters(client *redis.Client, app string, env string, cluster string, detail bool) (string, error) {

	var clusters models.ClusterList
	//mymap := make(map[string]interface{})

	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("clusters/%v-%v-clusterlist.yaml", app, env))
	if err != nil {
		return "", fmt.Errorf("Could not find a cluster definition for application: %v env: %v", app, env)
	}
	err = yaml.Unmarshal(yamlFile, &clusters)
	check(err)

	common.LogInfo.Printf("Description: %#v\n", clusters.Description)
	var output bytes.Buffer
	output.WriteString("{")
	isFirst := true
	for key, value := range clusters.Clusters {
		// If cluster specified only get info for that one cluster
		if cluster != "" && cluster != key {
			continue
		}
		fmt.Println("Checking Cluster", key)
		fmt.Println("Checking value", value)
		val, err := client.Get(fmt.Sprintf("%v-%v-%v", app, env, key)).Result()
		if err == redis.Nil {
			common.LogInfo.Printf("Cluster %v is free", key)

		} else if err != nil {
			common.LogError.Printf("Failed to query redis for key [%v-%v-%v]", app, env, key)
		} else {
			common.LogInfo.Printf("Cluster %v has a reservation valid for 0 seconds %v", key, val)

			out := map[string]interface{}{}
			json.Unmarshal([]byte(val), &out)

			out["application"] = app
			out["environment"] = env
			out["cluster"] = key
			keyExp := client.TTL(fmt.Sprintf("%v-%v-%v", app, env, key))
			resourceName := fmt.Sprintf("%v-%v-%v", app, env, key)
			out["ttl"] = keyExp.String()

			if detail == true {
				resources, err := KubeGetDeployments(resourceName)
				if err != nil {
					common.LogError.Println("Getting Deployments information form kubernetes for key ", key)
				} else {
					if len(resources) > 0 {
						out["resources"] = resources
					}
				}
			}
			outputJSON, _ := json.Marshal(out)
			if isFirst == true {
				isFirst = false
			} else {
				output.WriteString(",")
			}
			output.WriteString(fmt.Sprintf("\"%v\" :", resourceName))
			output.WriteString(string(outputJSON))
		}
	}
	output.WriteString("}")
	common.LogInfo.Printf("List Of clusters [%v]", output)

	//Make it pretty
	out := map[string]interface{}{}
	json.Unmarshal(output.Bytes(), &out)
	outputJSON, _ := json.MarshalIndent(out, "", " ")
	common.LogInfo.Println("Output :", output.String())

	return string(outputJSON), nil

}

func check(e error) {
	if e != nil {
		common.LogError.Println(e)
	}
}
