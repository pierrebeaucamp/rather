package rather

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

type Dare struct {
	OptionA   string
	OptionB   string
	AmountA   int
	AmountB   int
	Permalink string
}

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/get", random)
	http.HandleFunc("/save", save)
	http.HandleFunc("/submit", submit)
	http.HandleFunc("/question/", question)
}

func get(w http.ResponseWriter, r *http.Request) ([]Dare, error) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Dare").Ancestor(parentProject(c))

	var dares []Dare
	_, err := q.GetAll(c, &dares)
	return dares, err
}

func index(w http.ResponseWriter, r *http.Request) {
	dares, err := get(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := template.Must(template.New("list").ParseFiles("views/list.html"))
	err = t.ExecuteTemplate(w, "list", dares)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Each project gets assigend the same ancestor so we have faster reads
func parentProject(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Project", "parent-project", 0, nil)
}

func random(w http.ResponseWriter, r *http.Request) {
	dares, err := get(w, r)
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

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprint(w, string(j))
}

func save(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id, _, err := datastore.AllocateIDs(c, "Dare", parentProject(c), 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key := datastore.NewKey(c, "Dare", "", id, parentProject(c))
	url := fmt.Sprintf("/question/%d", id)

	d := Dare{
		OptionA:   r.FormValue("OptionA"),
		OptionB:   r.FormValue("OptionB"),
		AmountA:   0,
		AmountB:   0,
		Permalink: url,
	}

	_, err = datastore.Put(c, key, &d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func submit(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("submit").ParseFiles("views/submit.html"))
	err := t.ExecuteTemplate(w, "submit", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func question(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	d := new(Dare)

	id, err := strconv.ParseInt(strings.Split(r.URL.Path, "/")[2], 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key := datastore.NewKey(c, "Dare", "", id, parentProject(c))
	err = datastore.Get(c, key, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := template.Must(template.New("question").ParseFiles("views/question.html"))
	err = t.ExecuteTemplate(w, "question", d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
