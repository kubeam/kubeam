package server

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/kubeam/kubeam/common"
	git "gopkg.in/src-d/go-git.v4"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const signaturePrefix = "sha1="
const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

//GitHubData ...
type GitHubData struct {
	URL       string
	RepoName  string
	CommitID  string
	Timestamp string
	RepoID    float64
}

//GitHubHook ...
type GitHubHook struct {
	ID        string
	Event     string
	Signature string
	Payload   []byte
}

/*GetGithubRepos downloads Github repos*/
func GetGithubRepos(w http.ResponseWriter, r *http.Request) {
	// var dir, symlink string
	// defer cleanup(dir, symlink)
	hook, err := ParseGithubHook(r)

	if err != nil || hook == nil {
		w.Write([]byte(err.Error()))
	} else {
		data := make(map[string]interface{})
		err = json.Unmarshal(hook.Payload, &data)
		ghdata := ParseGithubPayload(data)
		ghdata.CloneAndSymlinkApp()
		ghdata.SaveGhData()
		w.Write([]byte("Success"))
	}
}

// ParseGithubHook reads a Hook from an incoming HTTP Request.
func ParseGithubHook(req *http.Request) (*GitHubHook, error) {
	var gh = NewHook(req)
	common.LogInfo.Printf("%##v", gh)

	if len(gh.Signature) == 0 {
		return nil, errors.New("No Github Signature")
	}
	if len(gh.Event) == 0 {
		return nil, errors.New("Empty Github Event")
	}
	if len(gh.ID) == 0 {
		return nil, errors.New("No Github Event ID")
	}

	gh.Payload, _ = ioutil.ReadAll(req.Body)
	secret, err := common.Config.GetString("github/secret", "")

	if err != nil || !gh.SignedBy([]byte(secret)) {
		return nil, errors.New("Unable to verify signature")
	}
	return gh, nil
}

/*NewHook inits a new GithubHook struct*/
func NewHook(req *http.Request) *GitHubHook {
	common.LogInfo.Printf("%##v", req.Header)
	return &GitHubHook{
		Signature: req.Header.Get("X-Hub-Signature"),
		Event:     req.Header.Get("X-GitHub-Event"),
		ID:        req.Header.Get("X-Github-Delivery"),
	}
}

// ParseGithubPayload reads a Hook from an incoming HTTP Request.
func ParseGithubPayload(m map[string]interface{}) *GitHubData {
	var repo map[string]interface{}
	var commits map[string]interface{}
	var ghdata = new(GitHubData)

	for k, v := range m {
		if strings.Compare(k, "repository") == 0 {
			repo = v.(map[string]interface{})
		}
		if strings.Compare(k, "head_commit") == 0 {
			commits = v.(map[string]interface{})
		}
	}
	for k, v := range repo {
		if strings.Compare(k, "clone_url") == 0 {
			ghdata.URL = v.(string)
		}
		if strings.Compare(k, "name") == 0 {
			ghdata.RepoName = v.(string)
		}
		if strings.Compare(k, "id") == 0 {
			ghdata.RepoID = v.(float64)
		}
	}
	for k, v := range commits {
		if strings.Compare(k, "id") == 0 {
			ghdata.CommitID = v.(string)
		}
		if strings.Compare(k, "timestamp") == 0 {
			ghdata.Timestamp = v.(string)
		}
	}
	return ghdata
}

/*CloneAndSymlinkApp downloads the github repo and symlinks it to
application/{repo}*/
func (ghdata *GitHubData) CloneAndSymlinkApp() (string, string) {
	dir, _ := os.Getwd()
	ndir := fmt.Sprintf("%s%s", dir, ghdata.RepoName)
	username, err := common.Config.GetString("github/username", "")
	token, err := common.Config.GetString("github/token", "")

	if err != nil {
		common.LogError.Println(err.Error())
	}

	_, err = git.PlainClone(ndir, false, &git.CloneOptions{
		URL: ghdata.URL,
		Auth: &githttp.BasicAuth{
			Username: username,
			Password: token,
		},
		Progress: os.Stdout,
	})
	if err != nil {
		common.LogError.Println(err.Error())
	}
	symlinkDest := fmt.Sprintf("%s%s", dir, ghdata.RepoName)
	symlinkSrc := fmt.Sprintf("/applications/%s", ghdata.RepoName)
	if err = os.Symlink(symlinkDest, symlinkSrc); err != nil {
		common.LogError.Println(err.Error())
	}
	return symlinkSrc, symlinkDest
}

// Remove symlinks and downloaded github repo to save space
func cleanup(dir, symlink string) {
	os.RemoveAll(dir)
	os.Remove(symlink)
}

/*SignedBy checks that the provided secret matches the hook Signature*/
func (h *GitHubHook) SignedBy(secret []byte) bool {
	if len(h.Signature) != signatureLength ||
		!strings.HasPrefix(h.Signature, signaturePrefix) {
		return false
	}

	messageMAC := make([]byte, 20)
	hex.Decode(messageMAC, []byte(h.Signature[5:]))
	expectedMAC := signBody(secret, h.Payload)
	return hmac.Equal(messageMAC, expectedMAC)
}

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return computed.Sum(nil)
}

/*SaveGhData ...*/
func (ghdata *GitHubData) SaveGhData() {
	db := GetDatabaseConnection()
	stmt, err := db.Prepare("REPLACE INTO gitrepos VALUES (?, ?, ?, ?)")
	common.ErrorHandler(err)
	_, err = stmt.Exec(ghdata.RepoID, ghdata.RepoName, ghdata.CommitID, ghdata.URL)
	common.ErrorHandler(err)
}

/*GetGhData ...*/
func GetGhData() ([](*GitHubData), error) {
	var ghdatalist [](*GitHubData)
	db := GetDatabaseConnection()
	stmt, err := db.Query("SELECT * from gitrepos")
	common.ErrorHandler(err)
	if err == nil {
		for stmt.Next() {
			var ghdata = new(GitHubData)
			if rows := stmt.Scan(&ghdata.RepoID, &ghdata.RepoName, &ghdata.CommitID, &ghdata.URL); rows != nil {
				return nil, rows
			}
			ghdatalist = append(ghdatalist, ghdata)
		}
		return ghdatalist, nil
	}
	return nil, err
}

/*LoadGitRepos ...*/
func LoadGitRepos(w http.ResponseWriter, r *http.Request) {
	var repolist []string
	ghdatalist, err := GetGhData()
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
