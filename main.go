package main // import "github.com/PierreZ/people-in-space-proxy-lametric"

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	// See http://golang.org/pkg/time/#Parse
	timeFormat = "2006-01-02"
)

var source string = "http://www.howmanypeopleareinspacerightnow.com/peopleinspace.json"

var peopleString = " people in space right now"

type PeopleInSpace struct {
	Number int `json:"number"`
	People []struct {
		Name           string `json:"name"`
		Biophoto       string `json:"biophoto"`
		Biophotowidth  int    `json:"biophotowidth"`
		Biophotoheight int    `json:"biophotoheight"`
		Country        string `json:"country"`
		Countryflag    string `json:"countryflag"`
		Launchdate     string `json:"launchdate"`
		Careerdays     int    `json:"careerdays"`
		Title          string `json:"title"`
		Location       string `json:"location"`
		Bio            string `json:"bio"`
		Biolink        string `json:"biolink"`
		Twitter        string `json:"twitter"`
	} `json:"people"`
}
type LaMetric struct {
	Frames []Frames `json:"frames"`
}
type Frames struct {
	Index int         `json:"index"`
	Text  string      `json:"text"`
	Icon  interface{} `json:"icon"`
}

func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServe(":3000", nil)
}

func getData() PeopleInSpace {

	res, err := http.Get(source)
	body, err := ioutil.ReadAll(res.Body)

	var data PeopleInSpace
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Panic("Error unmarshalling JSON")
	}
	return data
}

func foo(w http.ResponseWriter, r *http.Request) {

	data := getData()

	frames := []Frames{
		Frames{Index: 0, Text: strconv.Itoa(data.Number) + peopleString, Icon: "i1631"},
	}

var temp string
	for _, value := range data.People {

		then, err := time.Parse(timeFormat, value.Launchdate)
		if err != nil {
			log.Fatal("Error parsing the date")
		}
		duration := time.Since(then)
		temp = temp + value.Name + " since " + strconv.FormatFloat(math.Ceil(duration.Hours()/24.00), 'f', 0, 64) + " days - "

	}
		var people Frames
			people = Frames{Index: 1, Text: temp, Icon: "i1631"}
		frames = append(frames, people)

	var lametric LaMetric
	lametric.Frames = frames

	js, err := json.Marshal(lametric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
