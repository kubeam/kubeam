package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/kubeam/kubeam/common"
	"github.com/kubeam/kubeam/models"
	yaml "gopkg.in/yaml.v2"
)

/*ReserveFeatureCluster reserves a feature cluster*/
func ReserveFeatureCluster(w http.ResponseWriter, r *http.Request) {
	var clusters models.ClusterList
	vars := mux.Vars(r)
	res := make(map[string]string)

	app := vars["application"]
	env := vars["environment"]
	branch := vars["branch"]

	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("clusters/%v-%v-clusterlist.yaml", app, env))
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Could not find a cluster definition for application: %v env: %v", app, env)))
	}
	err = yaml.Unmarshal(yamlFile, &clusters)
	//
	// BUG/FIX: This should generate a error payload  with a json struct and return
	common.LogError.Println(err.Error())
	///check(err)

	ttl := time.Duration(604800 * time.Second)
	common.LogInfo.Println(app, env, branch, ttl)
	rediskey := fmt.Sprintf("%v-%v-%v", app, env, branch)

	cluster, err := FindAndReserveCluster(common.RedisClient, rediskey, clusters, ttl)

	if err != nil {
		common.LogError.Println(err.Error())
		res["Cluster"] = err.Error()
		res["TTL"] = "-1"
	} else {
		ttl, _ := common.RedisClient.TTL(rediskey).Result()
		res["Cluster"] = cluster
		res["TTL"] = fmt.Sprintf("%f", ttl.Seconds())
	}
	parsed, _ := json.Marshal(res)
	w.Write(parsed)
}

/*FindAndReserveCluster ...*/
func FindAndReserveCluster(client *redis.Client, rediskey string, clusters models.ClusterList, ttl time.Duration) (string, error) {
	cls := getAllocatedClusters(client)

	common.LogInfo.Println(cls)
	for key, value := range clusters.Clusters {
		common.LogDebug.Println("Checking Cluster", key)
		common.LogDebug.Println("Checking value", value)

		res, err := client.Get(rediskey).Result()
		if _, ok := cls[key]; !ok {
			common.LogDebug.Println(key)
			if err == redis.Nil {
				common.LogInfo.Println(fmt.Sprintf("Found available cluster [%v] for you.", key))
				featurettl, err := time.ParseDuration("604800s")
				ttl = time.Duration(featurettl * time.Second)

				err = ReserveCluster(client, rediskey, key, ttl)
				if err != nil {
					return "", err
				}
				return string(key), nil
			} else if err != nil {
				common.LogError.Println(err.Error())
				return "", fmt.Errorf("Failed to query redis for key [%v]", rediskey)
			} else {
				return res, nil
			}
		}
	}

	common.LogDebug.Printf("No clusters available for reservation")
	return "", fmt.Errorf("No clusters available for reservation for appliation")
}

/*ReserveCluster ...*/
func ReserveCluster(client *redis.Client, rediskey, cluster string, t time.Duration) error {
	if err := client.Set(rediskey, cluster, t).Err(); err != nil {
		return err
	}
	common.LogDebug.Printf("Created reservation for %v with TTL of %f", rediskey, t.Seconds())
	return nil
}

/*FreeFeatureCluster ...*/
func FreeFeatureCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rediskey := fmt.Sprintf("%s-%s-%s", vars["application"], vars["environment"], vars["branch"])
	_, err := common.RedisClient.Del(rediskey).Result()
	if err != nil {
		common.LogError.Println(err.Error())
	}
}

func getAllocatedClusters(client *redis.Client) map[string]bool {
	var clusters = make(map[string]bool)
	keys, err := client.Keys("*").Result()
	if err != nil {
		return nil
	}
	for _, val := range keys {
		cl, _ := client.Get(val).Result()
		clusters[cl] = true
	}
	return clusters
}

/*GetFeatureCluster ...*/
func GetFeatureCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res := make(map[string]string)
	rediskey := fmt.Sprintf("%s-%s-%s", vars["application"], vars["environment"], vars["branch"])

	if cluster, err := common.RedisClient.Get(rediskey).Result(); err != redis.Nil {
		ttl, _ := common.RedisClient.TTL(rediskey).Result()
		res["TTL"] = fmt.Sprintf("%f", ttl.Seconds())
		res["Cluster"] = cluster

		if response, err := json.Marshal(res); err == nil {
			w.Write(response)
		} else {
			w.Write([]byte(err.Error()))
		}
	} else {
		w.Write([]byte(err.Error()))
	}
}
