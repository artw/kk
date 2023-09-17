package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"
)

var dynamicClient dynamic.Interface

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
			panic(err.Error())
	}

	dynamicClient, err = dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Check the content type and convert if necessary
	contentType := r.Header.Get("Content-Type")
	var jsonBody []byte
	switch contentType {
	case "application/yaml":
		jsonBody, err = yaml.YAMLToJSON(body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error converting YAML to JSON: %v", err), http.StatusBadRequest)
			return
		}
	case "application/json":
		jsonBody = body
	default:
		http.Error(w, "Unsupported content type", http.StatusBadRequest)
		return
	}

	var unstructuredObj *unstructured.Unstructured
	if err := json.Unmarshal(jsonBody, &unstructuredObj); err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshaling body: %v", err), http.StatusBadRequest)
		return
	}

	gvk := unstructuredObj.GroupVersionKind()
	resource := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: fmt.Sprintf("%ss", strings.ToLower(gvk.Kind)), // naive pluralization
	}

	namespace := unstructuredObj.GetNamespace()

	switch r.Method {
	case "PUT":
		_, err = dynamicClient.Resource(resource).Namespace(namespace).Update(context.TODO(), unstructuredObj, v1.UpdateOptions{})
	case "POST":
		_, err = dynamicClient.Resource(resource).Namespace(namespace).Create(context.TODO(), unstructuredObj, v1.CreateOptions{})
	case "DELETE":
		resourceName := unstructuredObj.GetName()
		err = dynamicClient.Resource(resource).Namespace(namespace).Delete(context.TODO(), resourceName, v1.DeleteOptions{})
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("API server responded with: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Object processed")
}
