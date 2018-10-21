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


  modulesX := make(map[string]string)
  fmt.Println(modulesX)

  // 搞明白这里了， 如果是range一个map， 那么range的值就是map的key
  // 所以以太坊原来的RPC代码实际上是用 要调用的方法名来做法 key，然后value 是指向了方法的引用
  for name := range s.server.services {
    fmt.Println(name)
  }


  testMap := make(map[string]string)
  testMap["xx"] = "xx_val"
  testMap["yy"] = "yy_val"
  for key := range testMap {
    fmt.Println(key)
  }

}








