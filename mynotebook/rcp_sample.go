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


type serviceX struct {
  name            string
  typ             reflect.Type
  callbacks       callbacksX
  subscriptions   subscriptionsX
}


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
  serviceRe := make(map[string]*serviceX)
  serviceRe["testKey"] = servicex

  s := *RPCServiceX{x}
  modulesX := make(map[string]string)
  fmt.Println(modulesX)

  for name := range s.server.services {
    fmt.Println(name)
  }

}


func (s *Server) RegisterName(name string, rcvr interface{}) error {
  if s.services == nil {
    s.services = make(serviceRegistry)
  } 

  svc := new(service)
  svc.typ = reflect.TypeOf(rcvr)
  rcvrVal := reflect.ValueOf(rcvr)

  if name == "" {
    return fmt.Errorf("no service name for type %s", svc.typ.String())
  } 

  // 判断变量名的首位是否是大写， 需要变量名能否对外使用
  if !isExported(reflect.Indirect(rcvrVal).Type().Name()) {
    return fmt.Errorf("%s is not exported", reflect.Indirect(rcvrVal).Type().Name())
  }

  methods, subscriptions := suitableCallbacks(rcvrVal, svc.typ)

  // already a previous service register under given sname, merge methods/subscriptions
  if regsvc, present := s.services[name]; present {
    if len(methods) == 0 && len(subscriptions) == 0 {
      return fmt.Errorf("Service %T doesn't have any suitable methods/subscriptions to expose", rcvr)
    }
    for _, m := range methods {
      regsvc.callbacks[formatName(m.method.Name)] = m
    }
    for _, s := range subscriptions {
      regsvc.subscriptions[formatName(s.method.Name)] = s
    }
    return nil
  }

  svc.name = name
  svc.callbacks, svc.subscriptions = methods, subscriptions

  if len(svc.callbacks) == 0 && len(svc.subscriptions) == 0 {
    return fmt.Errorf("Service %T doesn't have any suitable methods/subscriptions to expose", rcvr)
  }

  s.services[svc.name] = svc
  return nil
  
}

























