![keva logo](https://github.com/donuts-are-good/keva/assets/96031819/89552a8e-949c-409e-aa55-e7f20cceaa69)
![donuts-are-good's followers](https://img.shields.io/github/followers/donuts-are-good?&color=555&style=for-the-badge&label=followers) ![donuts-are-good's stars](https://img.shields.io/github/stars/donuts-are-good?affiliations=OWNER%2CCOLLABORATOR&color=555&style=for-the-badge) ![donuts-are-good's visitors](https://komarev.com/ghpvc/?username=donuts-are-good&color=555555&style=for-the-badge&label=visitors)
# libkeva

libkeva is a library for [keva](https://github.com/donuts-are-good/keva), a key:value datastore with http json interface.

## demo:

```
package main

import (
	"fmt"
	"time"
	"github.com/donuts-are-good/libkeva"
)

func main() {
	store := libkeva.NewKeyValueStore("data.json", 5*time.Second)

	// Load initial data from file if it exists
	err := store.LoadFromFile("data.json")
	if err != nil {
		fmt.Println("Error loading data from file:", err)
		return
	}

	// Set some key-value pairs
	store.Set("key1", "value1")
	store.Set("key2", "value2")
	store.Set("key3", "value3")

	// Get values
	val1, _ := store.Get("key1")
	val2, _ := store.Get("key2")
	val3, _ := store.Get("key3")

	fmt.Println("Key1:", val1)
	fmt.Println("Key2:", val2)
	fmt.Println("Key3:", val3)

	// Delete a key
	store.Delete("key3")

	val3, exists := store.Get("key3")
	if !exists {
		fmt.Println("Key3 has been deleted")
	} else {
		fmt.Println("Key3:", val3)
	}

	// Let's sleep for a bit to allow any pending saves to complete.
	time.Sleep(5 * time.Second)

	// Reload the data from file to simulate a restart
	newStore := libkeva.NewKeyValueStore("data.json", 5*time.Second)
	err = newStore.LoadFromFile("data.json")
	if err != nil {
		fmt.Println("Error loading data from file:", err)
		return
	}

	// Check if our data is there
	val1, _ = newStore.Get("key1")
	val2, _ = newStore.Get("key2")

	fmt.Println("After reloading from file:")
	fmt.Println("Key1:", val1)
	fmt.Println("Key2:", val2)
}


```


### 1. store a key-value:



`curl -x post http://localhost:8080/store/demokey -h "content-type: application/json" -d '{"value": "demo value"}'`

**output:**

`ok`

### 2. retrieve a stored value:



`curl -x get http://localhost:8080/store/demokey`

**output:**



`"demo value"`

### 3. delete a stored key:



`curl -x delete http://localhost:8080/store/demokey`

**output:**

`ok`

### 4. try retrieving a deleted key:



`curl -x get http://localhost:8080/store/demokey`

**output:**


`key not found`

### 5. health check:



`curl -x get http://localhost:8080/health`

**output:**

`healthy`

## endpoints:

### 1. get /store/{key}

**description:** retrieve a value by the given key.

**parameters:**

- `key`: the key associated with the stored value.

**response:**

- `200 ok`: value retrieved successfully. it returns the stored value in json format.
- `404 not found`: key not found.

**example:**



`curl -x get http://localhost:8080/store/examplekey`

---

### 2. post /store/{key}

**description:** store a value associated with the given key.

**parameters:**

- `key`: the key to store the value with.

**request body:** json object containing the value to store.

- `value` (string): the value to be stored.

**response:**

- `201 created`: key-value set successfully.
- `400 bad request`: no value provided or bad request format.

**example:**



`curl -x post http://localhost:8080/store/examplekey -h "content-type: application/json" -d '{"value": "this is an example value"}'`

---

### 3. delete /store/{key}

**description:** delete the value associated with the given key.

**parameters:**

- `key`: the key of the value to delete.

**response:**

- `200 ok`: key deleted successfully.
- `404 not found`: key not found.

**example:**



`curl -x delete http://localhost:8080/store/examplekey`

---

### 4. get /health

**description:** health check endpoint.

**response:**

- `200 ok`: healthy.

**example:**



`curl -x get http://localhost:8080/health`

## errors:

the api uses conventional http response codes to indicate the success or failure of an api request.

- `200 ok`: the request was successful.
- `201 created`: the request was successful and a resource was created.
- `400 bad request`: the request could not be understood or was missing required parameters.
- `404 not found`: resource not found. this can be used when a specific key does not exist in the store.
- `405 method not allowed`: the http method used is not valid for the specific endpoint.

## license

mit license 2023 donuts-are-good, for more info see license.md
