package main 

import (
    "fmt"
    "strings"
    "errors"
    "encoding/json"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "gopkg.in/cavaliercoder/g2z.v3"
    "github.com/Jeffail/gabs"
)



func connect(mongoUrl string) (*mgo.Session, error) {
    g2z.LogInfof("Dialing server")
    session, err := mgo.Dial(mongoUrl)
    return session, err 
}


func queryDB(request *g2z.AgentRequest) (string, error) {
    // validate param count
    nbParams := len(request.Params)
    if nbParams != 3 {
        return "", errors.New(fmt.Sprintf("Invalid parameter count, expected 3, got %d", nbParams))
    }
    mongoUrl := request.Params[0]
    rawQuery := request.Params[1]
    rawWantedResult := request.Params[2]

    if mongoUrl == "" {
        mongoUrl = "mongodb://127.0.0.1/local"
    }

    g2z.LogInfof("Opening session with DB %s", mongoUrl)
    session, err := connect(mongoUrl) 
    if err != nil {
        g2z.LogErrorf("Could not connect to mongo URL %s: %s", mongoUrl, err)
        session.Close()
        return "", errors.New(fmt.Sprintf("%s", err))
    }

    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)
    db := session.DB("")
 
    result := bson.M{}
    query := bson.M{}

    // We try to parse the query as a JSON object
    if err := json.Unmarshal([]byte(rawQuery), &query); err != nil {
        g2z.LogErrorf("Could not parse query as valid JSON: %s", err)
        session.Close()
        return "", errors.New(fmt.Sprintf("%s", err))
    }
 
    // We now try to run the query
    if err := db.Run(query, &result); err != nil {
        g2z.LogErrorf("Could not run query %s: %s", query, err)
        session.Close()
        return "", errors.New(fmt.Sprintf("%s", err))
    }

   
    // Let's try to parse the result as JSON
    json, err := json.Marshal(result)
    if err != nil {
        g2z.LogErrorf("Could not parse query result as valid JSON: %s", err)
        session.Close()
        return "", errors.New(fmt.Sprintf("%s", err)) 
    } 
    g2z.LogDebugf("JSON result=%s", json)

    session.Close()

    if rawWantedResult == "" {
        return strings.TrimSpace(fmt.Sprintf("%s", json)), nil
    }

    jsonParsed, err := gabs.ParseJSON([]byte(fmt.Sprintf("%s", json)))

    if err != nil {
        g2z.LogErrorf("Gabs library could not parse JSONresult: %s", err)
        return "", errors.New(fmt.Sprintf("%s", err))
    }
    
    value := jsonParsed.Path(rawWantedResult).String()

    g2z.LogInfof("Wanted value=%s", value)
    return value, nil
}
