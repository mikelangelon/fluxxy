package main

import "flag"

var httpInterface = flag.String("http.interface", ":8080", "HTTP service interface as in <address>:<port>")

var mongoURL = flag.String("mongo.url", "localhost", "MongoDB connect string")
var mongoUser = flag.String("mongo.username", "", "MongoDB username")
var mongoPass = flag.String("mongo.password*", "", "MongoDB passord")
var stubsEnabled = flag.Bool("stubs.enable", true, "Enable test/stub endpoints on application")
