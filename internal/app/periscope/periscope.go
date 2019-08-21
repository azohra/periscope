package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/azohra/periscope/internal/pkg/periscope/models"
	"github.com/azohra/periscope/internal/pkg/periscope/services"
	"github.com/azohra/periscope/internal/pkg/periscope/tools"

	apiv1 "k8s.io/api/core/v1"
)

var count int

// Route main router
func Route() {
	count = 0
	fmt.Println("Launched periscope backend")
	http.HandleFunc("/allocate", AllocateHandler)
	http.HandleFunc("/ping", PingHandler)
	http.ListenAndServe(":7000", nil)
}

// PingHandler pong
func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pong")
}

// AllocateHandler increment box full of nuts
func AllocateHandler(w http.ResponseWriter, r *http.Request) {

	payloadStr := r.FormValue("payload")
	if payloadStr != "" {

		var payload models.Payload
		json.Unmarshal([]byte(payloadStr), &payload)

		fmt.Println(payload)
		fmt.Println(payloadStr)

		if len(payload.Images) > 0 {

			stateIDPrefix := "periscope-mock"
			stateIDStr := payload.StateID
			appName := fmt.Sprintf("%s-%s", stateIDPrefix, stateIDStr)

			containers := []apiv1.Container{}

			for _, elem := range payload.Images {
				container := apiv1.Container{
					Name:  fmt.Sprintf("%s%s", appName, tools.RandStr(10)),
					Image: elem.Image,
					Ports: []apiv1.ContainerPort{
						{
							Name:          "http",
							Protocol:      apiv1.ProtocolTCP,
							ContainerPort: int32(elem.Port),
						},
					},
				}
				containers = append(containers, container)
			}

			services.CreateDeployment(appName, "periscope", containers, "gcr-json-key")
			services.CreateService(appName, "periscope", int32(payload.SvcPort))
			services.CreateVirtualService(appName, "periscope", stateIDStr, payload.SvcPort, "v1")
			fmt.Fprintf(w, "Done.")

		} else {
			fmt.Println("Invalid payload")
			fmt.Fprintf(w, "Invalid Payload")
		}
	} else {
		fmt.Fprintf(w, "No Payload")
	}
}
