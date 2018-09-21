package handlers

import (
	"encoding/json"
	"github.com/kubeam/kubeam/common"
	"github.com/kubeam/kubeam/services"
	"net/http"
)

/*GetGithubRepos downloads Github repos*/
func GetGithubRepos(w http.ResponseWriter, r *http.Request) {
	// var dir, symlink string
	// defer cleanup(dir, symlink)
	hook, err := services.ParseGithubHook(r)

	if err != nil || hook == nil {
		w.Write([]byte(err.Error()))
	} else {
		data := make(map[string]interface{})
		err = json.Unmarshal(hook.Payload, &data)
		ghdata := services.ParseGithubPayload(data)
		ghdata.CloneAndSymlinkApp()
		ghdata.SaveGhData()
		w.Write([]byte("Success"))
	}
}

/*LoadGitRepos ...*/
func LoadGitRepos(w http.ResponseWriter, r *http.Request) {
	var repolist []string
	ghdatalist, err := services.GetGhData()
	if err != nil {
		common.LogError.Println(err.Error())
	} else {
		for _, ghdata := range ghdatalist {
			repolist = append(repolist, ghdata.RepoName)
			go ghdata.CloneAndSymlinkApp()
		}
	}
	res, err := json.Marshal(repolist)
	w.Write(res)
}
