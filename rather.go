package rather

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
)

type Dare struct {
	optionA string
	optionB string
	amountA int
	amountB int
}

/*
type User struct {
	phone string
}

type Choice struct {
	user   User
	dare   Dare
	choseA bool
}
*/

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/get", get)
	http.HandleFunc("/save", save)
}

func index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("index").ParseFiles("views/index.html"))
	err := t.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Dare").Ancestor(parentProject(c))

	var dares []Dare
	_, err := q.GetAll(c, &dares)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d := dares[rand.Intn(len(dares))]
	j, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(j))
}

func save(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	d := Dare{
		optionA: r.FormValue("OptionA"),
		optionB: r.FormValue("OptionB"),
		amountA: 0,
		amountB: 0,
	}

	key := datastore.NewIncompleteKey(c, "Dare", parentProject(c))
	_, err := datastore.Put(c, key, &d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// Each project gets assigend the same ancestor so we have faster reads
func parentProject(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Project", "parent-project", 0, nil)
}
