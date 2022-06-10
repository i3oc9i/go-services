// Package checkgrp maintains the group of handlers for health checking.
package checkgrp

import (
	"encoding/json"
	"net/http"
	"os"

	"go.uber.org/zap"
)

// Handlers manages the set of check endpoints.
type Handlers struct {
	Build string
	Log   *zap.SugaredLogger
}

// Readiness checks if the database is ready and if not will return a 500 status.
// Do not respond by just returning an error because further up in the call
// stack it will interpret that as a non-trusted error.
func (h Handlers) Readiness(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	statusCode := http.StatusOK

	if err := response(w, statusCode, data); err != nil {
		h.Log.Errorw("readiness", "ERROR", err)
	}

	h.Log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}

// Liveness returns simple status info if the service is alive.
// If the app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API.
// The Kubernetes environment variables have to be set within the Pod/Deployment manifest.
func (h Handlers) Liveness(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Status    string `json:"status,omitempty"`
		Build     string `json:"build,omitempty"`
		Node      string `json:"node,omitempty"`
		Namespace string `json:"namespace,omitempty"`
		Pod       string `json:"pod,omitempty"`
		PodIP     string `json:"podIP,omitempty"`
	}{
		Status:    "up",
		Build:     h.Build,
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
	}

	statusCode := http.StatusOK

	if err := response(w, statusCode, data); err != nil {
		h.Log.Errorw("liveness", "ERROR", err)
	}

	h.Log.Infow("liveness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}

// response writes well formed HTTP response.
func response(w http.ResponseWriter, statusCode int, data any) error {

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
