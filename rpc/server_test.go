// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rpc

import (
  "context"
  "encoding/json"
  "net"
  "reflect"
  "testing"
  "time"
  "fmt"
)

type Service struct{}

type Args struct {
  S string
}

func (s *Service) NoArgsRets() {
}

type Result struct {
  String string
  Int    int
  Args   *Args
}

func (s *Service) Echo(str string, i int, args *Args) Result {
  return Result{str, i, args}
}

func (s *Service) EchoWithCtx(ctx context.Context, str string, i int, args *Args) Result {
  return Result{str, i, args}
}

func (s *Service) Sleep(ctx context.Context, duration time.Duration) {
  select {
  case <-time.After(duration):
  case <-ctx.Done():
  }
}

func (s *Service) Rets() (string, error) {
  return "", nil
}

func (s *Service) InvalidRets1() (error, string) {
  return nil, ""
}

func (s *Service) InvalidRets2() (string, string) {
  return "", ""
}

func (s *Service) InvalidRets3() (string, string, error) {
  return "", "", nil
}

func (s *Service) Subscription(ctx context.Context) (*Subscription, error) {
  return nil, nil
}

/**
@test cmd:  go test -v server_test.go subscription.go types.go utils.go  server.go errors.go json.go -test.run TestServerRegisterName

 */
func TestServerRegisterName(t *testing.T) {
  server := NewServer()
  service := new(Service)

  if err := server.RegisterName("calc", service); err != nil {
    t.Fatalf("%v", err)
  }

  if len(server.services) != 2 {
    t.Fatalf("Expected 2 service entries, got %d", len(server.services))
  }

  svc, ok := server.services["calc"]
  if !ok {
    t.Fatalf("Expected service calc to be registered")
  }

  if len(svc.callbacks) != 5 {
    t.Errorf("Expected 5 callbacks for service 'calc', got %d", len(svc.callbacks))
  }

  if len(svc.subscriptions) != 1 {
    t.Errorf("Expected 1 subscription for service 'calc', got %d", len(svc.subscriptions))
  }
}

func TestServerRegisterNameMy(t *testing.T) {
  server := NewServer()
  fmt.Println(server)
  service := new(Service)

  // 注册服务名称
  if err := server.RegisterName("xxtest", service); err != nil {
    t.Fatalf("%v", err)
  }

  //fmt.Println(len(server.run))
  fmt.Println(server.services)
  fmt.Println(server.run)
  fmt.Println(server.codecs)

  if len(server.services) != 2 {
    t.Fatalf("Expected 2 service entries, got %d", len(server.services))
  }

  svc, ok := server.services["xxtest"]
  if !ok {
    t.Fatalf("Expected service xxtest to be registered")
  }

  fmt.Println(svc.callbacks)
  if len(svc.callbacks) != 5 {
    t.Errorf("Expected 5 callbacks for service 'xxtest', got %d", len(svc.callbacks))
  }

}

func testServerMethodExecution(t *testing.T, method string) {
  server := NewServer()
  service := new(Service)

  if err := server.RegisterName("test", service); err != nil {
    t.Fatalf("%v", err)
  }

  stringArg := "string arg"
  intArg := 1122
  argsArg := &Args{"abcde"}
  params := []interface{}{stringArg, intArg, argsArg}

  request := map[string]interface{}{
    "id":      12345,
    "method":  "test_" + method,
    "version": "2.0",
    "params":  params,
  }

  clientConn, serverConn := net.Pipe()
  defer clientConn.Close()

  go server.ServeCodec(NewJSONCodec(serverConn), OptionMethodInvocation)

  out := json.NewEncoder(clientConn)
  in := json.NewDecoder(clientConn)

  if err := out.Encode(request); err != nil {
    t.Fatal(err)
  }

  response := jsonSuccessResponse{Result: &Result{}}
  if err := in.Decode(&response); err != nil {
    t.Fatal(err)
  }

  if result, ok := response.Result.(*Result); ok {
    if result.String != stringArg {
      t.Errorf("expected %s, got : %s\n", stringArg, result.String)
    }
    if result.Int != intArg {
      t.Errorf("expected %d, got %d\n", intArg, result.Int)
    }
    if !reflect.DeepEqual(result.Args, argsArg) {
      t.Errorf("expected %v, got %v\n", argsArg, result)
    }
  } else {
    t.Fatalf("invalid response: expected *Result - got: %T", response.Result)
  }
}

func testServerMethodExecutionMy(t *testing.T, method string) {
  // 注册RPC服务
  server := NewServer()
  service := new(Service)

  if err := server.RegisterName("test", service); err != nil {
    t.Fatalf("%v", err)
  }

  // 定义 RPC 请求的参数
  stringArg := "string arg"
  intArg := 1122
  argsArg := &Args{"abcde"}
  params := []interface{}{ stringArg, intArg, argsArg }

  request := map[string]interface{}{
    "id":          12345,
    "method":      "test_" + method,
    "version":    "2.0",
    "params":      params,
  }

  // 设置 RPC 的网络连接,  net.Pipe() 就是一个网络读写的通道
  // 在本机测试，没有测试远程调用的网络， 这个只是 本地的 网络读写通道
  // Pipe创建一个内存中的同步、全双工网络连接。连接的两端都实现了Conn接口。一端的读取对应另一端的写入，直接将数据在两端之间作拷贝；没有内部缓冲。
  clientConn, serverConn := net.Pipe()
  defer clientConn.Close()

  go server.ServeCodec(NewJSONCodec(serverConn), OptionMethodInvocation)

  out := json.NewEncoder(clientConn)
  //int := json.NewEncoder(clientConn)
  fmt.Printf("out ------------------------------ %v\n", out)
  in := json.NewDecoder(clientConn)
  //out := json.NewDecoder(clientConn)
  fmt.Printf("in ------------------------------ %v\n", in)

  // out端 encode request 结构体之后， in 端可以直接获取
  if err := out.Encode(request); err != nil {
    t.Fatal(err)
  }

  response := jsonSuccessResponse{Result:  &Result{}}

  // 把 out端的 request 结构体进行 decode 并且获取
  if err := in.Decode(&response); err != nil {
    t.Fatal(err)
  }

  if result, ok := response.Result.(*Result); ok {
    if result.String != stringArg {
      t.Errorf("expected %s, got : %s\n", stringArg, result.String)
    }

    if result.Int != intArg {
      t.Errorf("expected %d, got : %d\n", intArg, result.Int)
    }

    if !reflect.DeepEqual(result.Args, argsArg) {
      t.Errorf("expected %v, got %v\n", argsArg, result)
    }

    fmt.Printf("response --------------------------------- %v\n", response)
    fmt.Printf("result --------------------------------- %v\n", result)
  } else {
    t.Fatalf("invalid response: expected *Result - got: %T", response.Result)
  }
}

func TestServerMethodExecution(t *testing.T) {
  testServerMethodExecution(t, "echo")
}


func TestServerMethodExecutionMy(t *testing.T) {
  testServerMethodExecutionMy(t, "echo")
}

func TestServerMethodWithCtx(t *testing.T) {
  testServerMethodExecution(t, "echoWithCtx")
}














