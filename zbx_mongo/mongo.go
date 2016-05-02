package main 

import (
    //"runtime"
    "fmt"
    "errors"
    "time"
    "strings"
    "encoding/json"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "gopkg.in/cavaliercoder/g2z.v3"
    "github.com/mattn/go-scan"
)


func validateParams(request *g2z.AgentRequest) (string, string, string, string, error){
    // validate param count
    nbParams := len(request.Params)
    if nbParams != 4 {
        return "", "", "", "", errors.New(fmt.Sprintf("Invalid parameter count, expected 4, got %d", nbParams))
    }
    mongoUrl := request.Params[0]
    command := request.Params[1]
    rawQuery := request.Params[2]
    rawWantedResult := request.Params[3]

    if mongoUrl == "" {
        mongoUrl = "mongodb://127.0.0.1/local"
    }
    if command == "" {
        return "", "", "", "", errors.New("Command cannot be null or empty (Parameter #2)")
    } 
    if rawWantedResult == "" {
        rawWantedResult = "/"
    }
    return mongoUrl, command, rawQuery, rawWantedResult, nil

}

func connect(mongoUrl string) (*mgo.Session, error) {
    timeout := g2z.Timeout - 1
    g2z.LogInfof("[zbx-mongo] Dialing server, timeout %ds", timeout)
    maxWait := time.Duration(timeout) * time.Second
    session, err := mgo.DialWithTimeout(mongoUrl, maxWait)
    if session != nil {
        session.SetSocketTimeout(maxWait)
        session.SetSyncTimeout(maxWait)
    }
    return session, err
}

// We have to build an ordered Mongo Query, where
// the command (find, dbStats...) has to be the first element,
// due to Json unmarshaller not respecting order
func prepareQuery(command string, unorderedQuery interface{}) (bson.D, error) {
    mQuery := unorderedQuery.(map[string]interface{})
    orderedQuery := make(bson.D, 1)
    commandFound := false
    
    // We have to build an ordered Mongo Query, where 
    // the command (find, dbStats...) has to be the first element, from an
    // unordered JSON document.
    for k, v := range mQuery {
        if strings.EqualFold(command, k) {
            orderedQuery[0] = bson.DocElem{k, v}
            commandFound = true
        } else {
           orderedQuery = append(orderedQuery, bson.DocElem{k, v})
        }
    }
    if commandFound == false {
        return nil, errors.New(fmt.Sprintf("Didn't find command '%s' in query", command))
    }
    return orderedQuery, nil
}

func queryDB(request *g2z.AgentRequest) (string, error) {
    mongoUrl, command, rawQuery, rawWantedResult, err := validateParams(request)
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
    //session.SetMode(mgo.Monotonic, true)
    db := session.DB("")
 
    var unorderedQuery interface{}
    result := bson.M{}

    // We try to parse the query as a JSON object
    if err := json.Unmarshal([]byte(rawQuery), &unorderedQuery); err != nil {
        g2z.LogErrorf("[zbx-mongo] Could not parse query as valid JSON: %s", err)
        session.Close()
        return "", errors.New(fmt.Sprintf("%s", err))
    }

    orderedQuery, err := prepareQuery(command, unorderedQuery)
    if err != nil {
        g2z.LogErrorf("[zbx-mongo] Could prepare query: %s", err)
        session.Close()
        return "", errors.New(fmt.Sprintf("%s", err)) 
    }
 
    // We now try to run the query
    if err := db.Run(orderedQuery, &result); err != nil {
        g2z.LogErrorf("[zbx-mongo] Could not run query %s: %s", orderedQuery, err)
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

    // If it's a simple string we want to remove the enclosing quotes
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
    //mongoUrl, rawParameters, rawWantedResult, err := validateParams(request)
    
    //if err != nil {
    //    g2z.LogErrorf("Error while validating parameters: %s", err)
    //    return nil, errors.New(fmt.Sprintf("%s", err))
    //}

    d := make(g2z.DiscoveryData, 5)
    return d, nil
}

