package middlefunc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	mod "github.com/ownperception/TechP_DB_Forum/models"
)

func Jsonparams(r *http.Request) (map[string]string, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close() // важный пункт!
	Check(err)
	arr := map[string]string{}
	err = json.Unmarshal(body, &arr)
	return arr, err
}
func GetJsonVote(r *http.Request) (mod.JsonVote, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close() // важный пункт!
	Check(err)
	vote := mod.JsonVote{}
	err = json.Unmarshal(body, &vote)
	return vote, err
}

func GetJsonArrayPost(r *http.Request) ([]mod.JsonPost, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close() // важный пункт!
	Check(err)
	posts := []mod.JsonPost{}
	err = json.Unmarshal(body, &posts)
	return posts, err
}

func JsonArrayparams(r *http.Request) ([]map[string]string, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close() // важный пункт!
	Check(err)
	arr := []map[string]string{}
	err = json.Unmarshal(body, &arr)
	return arr, err
}

func ParsUrl(r *http.Request, params *map[string]string) {
	for key, _ := range *params {
		param, ok := r.URL.Query()[key]
		if ok && len(param[0]) >= 1 {
			(*params)[key] = param[0]
		}
	}
}
