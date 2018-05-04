package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"gopkg.in/yaml.v2"
)

type ClusterList struct {
	Description string
	Clusters    map[string]map[string]string
}

var redisClient *redis.Client

func NewDBClient() *redis.Client {
	redisHost, err := config.GetString("redis/host", "localhost")
	redisPort, err := config.GetInt("redis/port", 6379)
	redisPassword, err := config.GetString("redis/password", "")
	LogInfo.Println(redisHost)
	LogInfo.Println(redisPort)
	LogInfo.Println(redisPassword)

	// BUG/FIX: should be using sentinel. We already have it running might as whell use it.
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", redisHost, redisPort),
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return (client)
}

func DBClientReserveCluster(client *redis.Client, app string, env string, key string, val []byte, t time.Duration) error {
	err := client.Set(fmt.Sprintf("%v-%v-%v", app, env, key), val, t).Err()
	//first slice then convert to string (string is a read-only slice of bytes)
	LogInfo.Println(fmt.Sprintf("%v-%v-%v [%v]", app, env, key, string(val[:])))
	if err != nil {
		return err
	} else {
		LogInfo.Println(fmt.Sprintf("Created reservation for %v-%v-%v [%v] with TTL of %v", app, env, key, string(val[:]), 300))
	}
	return nil
}

func DBClientFindAndReserve(client *redis.Client, app string, env string, ttl time.Duration) (string, error) {

	var clusters ClusterList

	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("clusters/%v-%v-clusterlist.yaml", app, env))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not find a cluster definition for application: %v env: %v", app, env))
	}
	err = yaml.Unmarshal(yamlFile, &clusters)
	check(err)

	LogInfo.Printf("Description: %#v\n", clusters.Description)
	for key, value := range clusters.Clusters {
		fmt.Println("Checking Cluster", key)
		fmt.Println("Checking value", value)
		val, err := client.Get(fmt.Sprintf("%v-%v-%v", app, env, key)).Result()
		if err == redis.Nil {
			LogInfo.Println(fmt.Sprintf("Found available cluster [%v] for you.", key))
			decodedValue, _ := json.Marshal(value)
			ttl, err := time.ParseDuration("1500s")
			if err != nil {
				default_ttl, _ := strconv.Atoi(value["default_ttl"])
				default_ttl_parsed := time.Duration(default_ttl)
				ttl = time.Duration(default_ttl_parsed * time.Second)
			}
			err = DBClientReserveCluster(client, app, env, string(key), decodedValue, ttl)
			if err != nil {
				return "", err
			} else {
				return string(key), nil
			}
		} else if err != nil {
			LogError.Println(fmt.Sprintf("Failed to query redis for key [%v-%v-%v]", app, env, key))
		} else {
			LogInfo.Println(fmt.Printf("Cluster %v has a reservation valid for 0 seconds %v", key, val))
		}
	}
	LogInfo.Println(fmt.Sprintf("No clusters available for reservation for appliation %v environment %v", app, env))
	return "", errors.New(fmt.Sprintf("No clusters available for reservation for appliation %v environment %v", app, env))

}

func DBGetClusterReservation(client *redis.Client, app string, env string, cluster string) (string, error) {

	var clusters ClusterList

	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("clusters/%v-%v-clusterlist.yaml", app, env))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not find a cluster definition for application: %v env: %v", app, env))
	}
	err = yaml.Unmarshal(yamlFile, &clusters)
	check(err)

	LogInfo.Printf("Description: %#v\n", clusters.Description)
	//gotReservation := false
	var output bytes.Buffer
	output.WriteString("{")
	val, err := client.Get(fmt.Sprintf("%v-%v-%v", app, env, cluster)).Result()
	if err == redis.Nil {
		LogInfo.Println(fmt.Sprintf("Cluster %v is free", cluster))
	} else if err != nil {
		LogError.Println(fmt.Sprintf("Failed to query redis for key [%v-%v-%v]", app, env, cluster))
	} else {
		LogInfo.Println(fmt.Printf("Cluster %v has a reservation valid for 0 seconds %v", cluster, val))

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
	LogInfo.Printf("List Of clusters [%v]", output)
	return output.String(), err

}

func DBClientGetSingleCluster(client *redis.Client, app string, env string, cluster string) (string, error) {
	ret, err := DBClientListClusters(client, app, env, cluster, false)
	return ret, err
}
func DBClientGetSingleClusterDetail(client *redis.Client, app string, env string, cluster string) (string, error) {
	ret, err := DBClientListClusters(client, app, env, cluster, true)
	return ret, err
}

func DBClientGetAllClusters(client *redis.Client, app string, env string) (string, error) {
	ret, err := DBClientListClusters(client, app, env, "", false)
	return ret, err
}
func DBClientGetAllClustersDetail(client *redis.Client, app string, env string) (string, error) {
	ret, err := DBClientListClusters(client, app, env, "", true)
	return ret, err
}

func DBClientListClusters(client *redis.Client, app string, env string, cluster string, detail bool) (string, error) {

	var clusters ClusterList
	//mymap := make(map[string]interface{})

	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("clusters/%v-%v-clusterlist.yaml", app, env))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not find a cluster definition for application: %v env: %v", app, env))
	}
	err = yaml.Unmarshal(yamlFile, &clusters)
	check(err)

	LogInfo.Printf("Description: %#v\n", clusters.Description)
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
			LogInfo.Println(fmt.Sprintf("Cluster %v is free", key))

		} else if err != nil {
			LogError.Println(fmt.Sprintf("Failed to query redis for key [%v-%v-%v]", app, env, key))
		} else {
			LogInfo.Println(fmt.Printf("Cluster %v has a reservation valid for 0 seconds %v", key, val))

			out := map[string]interface{}{}
			json.Unmarshal([]byte(val), &out)

			out["application"] = app
			out["environment"] = env
			out["cluster"] = key
			keyExp := client.TTL(fmt.Sprintf("%v-%v-%v", app, env, key))
			resourceName := fmt.Sprintf("%v-%v-c%v", app, env, key)
			out["ttl"] = keyExp.String()

			if detail == true {
				resources, err := KubeGetDeployments(resourceName)
				if err != nil {
					LogError.Println("Getting Deployments information form kubernetes for key ", key)
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
	LogInfo.Printf("List Of clusters [%v]", output)

	//Make it pretty
	out := map[string]interface{}{}
	json.Unmarshal(output.Bytes(), &out)
	outputJSON, _ := json.MarshalIndent(out, "", " ")
	LogInfo.Println("Output :", output.String())

	return string(outputJSON), nil

}

func check(e error) {
	if e != nil {
		LogError.Println(e)
	}
}
