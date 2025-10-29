# 🍕 PizzaLand Protocol Buffers

> **This directory defines the gRPC API contracts for the PizzaLand backend.**  
> It contains `.proto` schemas, a Makefile for Go code generation, and utilities for message validation.

---

## 📜 Overview

The **PizzaLand API** defines all message structures and service methods for managing pizzas and categories in the system.

This layer is completely independent of business logic — it’s used to generate Go source code for the backend and can also be reused by other services (e.g., frontend gateways, microservices, tests).

---

## 🧩 Service Definition

The `PizzaLand` gRPC service provides a complete CRUD interface for pizzas and their categories:

```protobuf
service PizzaLand {
  rpc Save(SaveRequest) returns (SaveResponse);                       // Save pizza
  rpc Get(GetRequest) returns (GetResponse);                           // Get pizza
  rpc List(ListRequest) returns (ListResponse);                        // List pizzas
  rpc Update(UpdateRequest) returns (UpdateResponse);                  // Update pizza
  rpc Remove(RemoveRequest) returns (RemoveResponse);                  // Remove pizza
  rpc SaveCategory(SaveCategoryRequest) returns (SaveCategoryResponse); // Save category
  rpc GetCategory(GetCategoryRequest) returns (GetCategoryResponse);     // Get category list
  rpc UpdateCategory(UpdateCategoryRequest) returns (UpdateCategoryResponse); // Update category
  rpc RemoveCategory(RemoveCategoryRequest) returns (RemoveCategoryResponse); // Remove category
}
```

Each RPC call uses validated request/response messages defined in `pizzaland.proto`.
The API enforces **strong typing**, **field validation**, and **Google-style annotations**.

---

## 🗂️ Directory Structure

```
api/
├── proto/
│   └── pizzaland/
│       └── pizzaland.proto     # Main proto definitions
├── generated/
│   └── go/                     # Generated Go gRPC and message code
├── Makefile                    # Protobuf build automation
└── README.md                   # You're here
```

---

## 🛠️ Protobuf Code Generation

All code generation tasks are automated using the provided **Makefile**.

### 🔧 Install Dependencies

Install or update the required Go protobuf tools:

```bash
  make install
```

This will install:

* [`protoc-gen-go`](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go)
* [`protoc-gen-go-grpc`](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc)
* [`protoc-gen-validate`](https://github.com/envoyproxy/protoc-gen-validate)

### 🧬 Generate Go Code

To generate Go protobuf and gRPC files:

```bash
  make gen
```

To generate protobuf + gRPC + validation code (recommended):

```bash
  make gen-validate
```

To perform all installation and generation steps (initial generating):

```bash
  make genall
```

> ⚠️ Always rerun `make genall` after editing any `.proto` file.

---

## 📦 Output Example

After running generation, Go code is placed inside `generated/go/`:

```
generated/go/
├── pizzaland.pb.go
├── pizzaland_grpc.pb.go
└── pizzaland_validate.pb.go
```

These contain:

* **pb.go** → data structures and message types
* **grpc.pb.go** → gRPC service interfaces and client/server stubs
* **validate.pb.go** → input validation rules from `protoc-gen-validate`

---

## 🧠 Proto Design Highlights

* **Validation Rules** with `protoc-gen-validate`
* **Field Annotations** using `google.api.field_behavior`
* **Use of OneOf** for flexible identifiers (`pizza_id` or `pizza_name`)
* **Enums** for strong type guarantees (`TypeDough`)
* **Nested Messages** for structured pizza and category data

Example excerpt:

```protobuf
// Next two files you need to download and place them in api/proto/...
import "third_party/googleapis/google/api/field_behavior.proto"; // download file 
import "validate/validate.proto"; // download file
import "google/protobuf/wrappers.proto";

message PizzaProperties {
  string name = 3 [
    (validate.rules).string = {min_len: 3, max_len: 50},
    (google.api.field_behavior) = REQUIRED
  ];
  float price = 6 [
    (validate.rules).float.gte = 109,
    (google.api.field_behavior) = REQUIRED
  ];
  TypeDough type_dough = 5 [
    (validate.rules).enum = {defined_only: true, not_in: [0]},
    (google.api.field_behavior) = REQUIRED
  ];
}
```

---

## 🧾 License

This project is licensed under the [MIT License](LICENSE).

---

### 🍕 *Proto-first. gRPC-fast. Pizza-ready.*
