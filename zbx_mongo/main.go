package main 


import (
    "runtime"
    "gopkg.in/cavaliercoder/g2z.v3"
)

var VERSION = "0.3"

func main() {
    panic("THIS_SHOULD_NEVER_HAPPEN")
}


func init() {
    //runtime.GOMAXPROCS(1)
    //runtime.LockOSThread()
    g2z.LogInfof("[zbx-mongo] Version %s loaded, using runtime %s.", VERSION, runtime.Version())

    g2z.RegisterStringItem("mongo.run", "", queryDB)
    g2z.RegisterDiscoveryItem("mongo.discover", "", discover)
}

