package main

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "net/http"
    "encoding/json"
)

type RIReportRequests struct {
    Id string `json:"_id" bson:"_id"`
    Chain string `json:"chain" bson:"chain"`
    ChainId string `json:"chainId" bson:"chainId"`
    Name string `json:"name" bson:"name"`
    HistoryDays int `json:"historyDays" bson:"historyDays"`
    AutoRefresh bool `json:"autoRefresh" bson:"autoRefresh"`
}
var colection *mgo.Collection

func handler(w http.ResponseWriter, r *http.Request) {
    var result []interface{}

    colection.Find(bson.M{}).All(&result)

    js, err := json.Marshal(result)
    if err != nil {
        //log.Fatal(err)
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

func main() {
    session, err := mgo.Dial("##############")
    if err != nil {
        panic(err)
    }
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)

    colection = session.DB("############").C("###########")

    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}