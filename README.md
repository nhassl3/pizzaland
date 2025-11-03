![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nhassl3/pizzaland) ![GitHub contributors](https://img.shields.io/github/contributors/nhassl3/pizzaland)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/nhassl3/pizzaland) ![GitHub last commit](https://img.shields.io/github/last-commit/nhassl3/pizzaland)

# 🍕 PizzaLand Backend — gRPC API

> **PizzaLand** is a modern backend for a pizza management system written in **Go**, using **gRPC** and **Protocol Buffers** for high-performance, type-safe communication.

This module represents the **API layer** of the project — it contains:
- All `.proto` definitions (service contracts, messages)
- The Makefile for generating Go stubs
- Code generation utilities for validation, protobuf, and gRPC

---

## 🚀 Overview

The **PizzaLand API** defines services and RPC methods for managing pizzas and categories in the system.  
It is designed to be clean, scalable, and ready for production.

### 🧩 Service Definition

```markdown
service PizzaLand {
    rpc Save(SaveRequest) returns (SaveResponse); // Save pizza procedure
    rpc Get(GetRequest) returns (GetResponse); // Get pizza procedure
    rpc List(ListRequest) returns (ListResponse); // Get list of the pizza procedure
    rpc Update(UpdateRequest) returns (UpdateResponse); // Update pizza properties or price procedure
    rpc Remove(RemoveRequest) returns (RemoveResponse); // Remove pizza from system procedure
    rpc SaveCategory(SaveCategoryRequest) returns (SaveCategoryResponse); // Save category for pizza on the system procedure
    rpc GetCategory(GetCategoryRequest) returns (GetCategoryResponse); // Get category object
    rpc GetCategoryList(GetCategoryListRequest) returns (GetCategoryListResponse); // Get list of the category procedure
    rpc UpdateCategory(UpdateCategoryRequest) returns (UpdateCategoryResponse); // Update category properties procedure
    rpc RemoveCategory(RemoveCategoryRequest) returns (RemoveCategoryResponse); // Remove category from the system procedure
}
```

---

## 🛠️ Project Structure

```
pizzaland/
├── proto/
│   └── pizzaland/
│       └── pizzaland.proto      # Main service and message definitions
├── generated/
│   └── go/                      # Generated Go source files
├── Makefile                     # Code generation automation
└── README.md                    # You're here 🍕
```

---

## ⚙️ Code Generation

All protobuf and gRPC Go stubs are generated via **Makefile** commands.

### 🧰 Prerequisites

Install the required tools (if not yet installed):

```bash
  make install
```

This will install:

* `protoc-gen-go`
* `protoc-gen-go-grpc`
* `protoc-gen-validate`

### 🔄 Generate bin and run

Generate bin

```bash
  make build
```

Run service:

```bash
  make run
```

Run tests or bench

```bash
  make test
```

```bash
  make bench
```

---

## 🧪 Example Generation Output

After running `make genall`, your `generated/go/` directory will contain:

```
generated/go/
├── pizzaland.pb.go
├── pizzaland_grpc.pb.go
└── pizzaland_validate.pb.go
```

These files include:

* Go data structures for messages (`*.pb.go`)
* Server and client stubs for gRPC (`*_grpc.pb.go`)
* Validation logic from `protoc-gen-validate` (`*_validate.pb.go`)

---

## 📦 Integration

You can import the generated Go code into your backend services:

```go
package client

import (
	"context"
	"fmt"
	"log"
	"net"

	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
)

func ExampleClient() {
	// Example usage
	conn, err := net.Dial("tcp", "localhost:44044")
	if err != nil {
		log.Fatal(err)
    }
	
	client := pizzalndv1.NewPizzaLandClient(conn)
	res, err := client.Get(context.Background(), &pizzalndv1.GetRequest{
		Identifier: &pizzalndv1.GetRequest_PizzaId{
			PizzaId: 1,
        },
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Pizza:", res.Pizza)
}

// Somewhere in main file (_ = main for example)
func _() {
	ExampleClient()
}
```

---

## 💡 Tips

* Keep `.proto` files small and modular.
* Always re-run `make genall` after editing `.proto` files.
* Use `buf` or `protolint` to ensure protobuf consistency and style.
* Consider versioning API definitions (`pizzaland/v1/`, `pizzaland/v2/`, etc.) as your project grows.

---

## 🧑‍🍳 Author & Links

**Author (frontend):** [LilBKb](https://github.com/LilBKb)\
**Author (backend):** [nhassl3](https://github.com/nhassl3)\
**Main Project:** [github.com/LilBKb/pizza](https://github.com/LilBKb/pizza)

---

## 🧾 License

This project is licensed under the [MIT License](LICENSE).