package main 

import (
    "fmt"
    "errors"
    "time"
    "encoding/json"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "gopkg.in/cavaliercoder/g2z.v3"
    "github.com/mattn/go-scan"
)


func validateParams(request *g2z.AgentRequest) (string, string, string, error){
    // validate param count
    nbParams := len(request.Params)
    if nbParams != 3 {
        return "", "", "", errors.New(fmt.Sprintf("Invalid parameter count, expected 3, got %d", nbParams))
    }
    mongoUrl := request.Params[0]
    rawQuery := request.Params[1]
    rawWantedResult := request.Params[2]

    if mongoUrl == "" {
        mongoUrl = "mongodb://127.0.0.1/local"
    }
    if rawQuery == "" {
        return "", "", "", errors.New("Bson Query cannot be null or empty (Parameter #2)")
    } 
    if rawWantedResult == "" {
        rawWantedResult = "/"
    }
    return mongoUrl, rawQuery, rawWantedResult, nil

}

func connect(mongoUrl string) (*mgo.Session, error) {
    g2z.LogInfof("Dialing server")
    maxWait := time.Duration(6 * time.Second)
    session, err := mgo.DialWithTimeout(mongoUrl, maxWait)
    return session, err
}

func queryDB(request *g2z.AgentRequest) (string, error) {

    mongoUrl, rawQuery, rawWantedResult, err := validateParams(request)
    if err != nil {
        g2z.LogErrorf("[zbx-mongo] Error while validating parameters: %s", err)
        return "", errors.New(fmt.Sprintf("%s", err))
    }

    g2z.LogInfof("[zbx-mongo] Opening session with DB %s", mongoUrl)
    session, err := connect(mongoUrl) 
    if err != nil {
        g2z.LogErrorf("[zbx-mongo] Could not connect to mongo URL %s: %s", mongoUrl, err)
        if session != nil {
            session.Close()
        }
        return "", errors.New(fmt.Sprintf("%s", err))
    }

    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)
    db := session.DB("")
 
    result := bson.M{}
    query := bson.M{}

    // We try to parse the query as a JSON object
    if err := json.Unmarshal([]byte(rawQuery), &query); err != nil {
        g2z.LogErrorf("[zbx-mongo] Could not parse query as valid JSON: %s", err)
        session.Close()
        return "", errors.New(fmt.Sprintf("%s", err))
    }
 
    // We now try to run the query
    if err := db.Run(query, &result); err != nil {
        g2z.LogErrorf("[zbx-mongo] Could not run query %s: %s", query, err)
        session.Close()
        return "", errors.New(fmt.Sprintf("%s", err))
    }

    session.Close()
    
    var typedValue interface{}
    var value string 
    err = scan.ScanTree(result, rawWantedResult, &typedValue)
    if err != nil {
        g2z.LogErrorf("[zbx-mongo] Could not search for '%s' in result: %s", rawWantedResult, err)
        return "", nil 
    }

    // It it's a simple string we want to remove the enclosing quotes
    switch typedValue.(type) {
    case string:
        value = fmt.Sprintf("%s", typedValue)
    default:
        jsonValue, err := json.Marshal(typedValue)
        if err != nil {
            g2z.LogErrorf("Could not parse result as valid JSON: %s", err)
        }
        value = fmt.Sprintf("%s", jsonValue)
    }
        
    g2z.LogInfof("[zbx-mongo] Wanted value=%s", value)
    return value, nil
}


func discover(request *g2z.AgentRequest) (g2z.DiscoveryData, error) {
    //mongoUrl, rawQuery, rawWantedResult, err := validateParams(request)
    
    //if err != nil {
    //    g2z.LogErrorf("Error while validating parameters: %s", err)
    //    return nil, errors.New(fmt.Sprintf("%s", err))
    //}

    d := make(g2z.DiscoveryData, 5)
    return d, nil
}

