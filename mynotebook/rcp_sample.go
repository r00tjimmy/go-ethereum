package main

import (
  "reflect"
  "sync"
  "gopkg.in/fatih/set.v0"
  "fmt"
)


type callbackX struct {
  rcvr            reflect.Value
  method          reflect.Method
  argTypes        []reflect.Type
  hasCtx          bool
  errPos          int
  isSubscribe     bool
}


// include:  callbacksX
type serviceX struct {
  name            string
  typ             reflect.Type
  callbacks       callbacksX
  subscriptions   subscriptionsX
}


// include: serviceRegistry
type ServerX struct {
  services      serviceRegistry

  run           int32
  codecsMu      sync.Mutex
  codecs        *set.Set
}


type serviceRegistry map[string]*serviceX
type callbacksX map[string]*callbackX
type subscriptionsX map[string]*callbackX


type RPCServiceX struct {
  server       *ServerX
}

func main() {
  servicex := &serviceX{name: "xxx"}
  fmt.Println(servicex) 
  fmt.Println(servicex.name)

  servicex2 := serviceX{name: "yyy"}
  fmt.Println(servicex2) 
  fmt.Println(servicex2.name)

  serviceRe := make(map[string]*serviceX)
  serviceRe["testKey"] = servicex

  ServerXTmp := &ServerX{ services: serviceRe }

  // ServerXTmp is the memory address
  s := &RPCServiceX{ server: ServerXTmp }
  fmt.Println(s)


  // modulesX := make(map[string]string)
  // fmt.Println(modulesX)

  // for name := range s.server.services {
  //   fmt.Println(name)
  // }

}








