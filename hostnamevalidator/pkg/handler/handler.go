package handler

// credit to https://github.com/giantswarm/grumpy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/rs/zerolog/log"
	"k8s.io/api/admission/v1beta1"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VerifyHandler listen to admission requests and serve responses
type VerifyHandler struct {
}

const baseDomain = "k8s.zach"

func (gs *VerifyHandler) Serve(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		log.Error().Msg("Empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}
	log.Info().Msg("Received request")

	if r.URL.Path != "/validate" {
		log.Error().Msg("No validate")
		http.Error(w, "no validate", http.StatusBadRequest)
		return
	}

	// unmarshal the admission request
	arRequest := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(body, &arRequest); err != nil {
		log.Error().Msg("Incorrect body")
		http.Error(w, "incorrect body", http.StatusBadRequest)
	}

	raw := arRequest.Request.Object.Raw
	// unmarshal the ingress
	ingress := networkv1.Ingress{}
	if err := json.Unmarshal(raw, &ingress); err != nil {
		log.Error().Msg("error deserializing ingress")
		return
	}

	// loop over the rule definitions for each host and validate the hostname
	for _, in := range ingress.Spec.Rules {
		valid, _ := evalHostname(in.Host, arRequest.Request.Namespace)
		if !valid {
			arResponse := v1beta1.AdmissionReview{
				Response: &v1beta1.AdmissionResponse{
					Allowed: false,
					Result: &metav1.Status{
						Message: "Hostname does match hostname requirements (hostname.namespace.k8s.zach)!",
					},
				},
			}
			resp, err := json.Marshal(arResponse)
			if err != nil {
				log.Err(err).Msg("Can't encode response")
				http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
			}
			log.Info().Msg("Ready to write reponse ...")
			if _, err := w.Write(resp); err != nil {
				log.Err(err).Msg("Can't write response")
				http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
			}
		}
	}

}

func evalHostname(hostname string, namespace string) (bool, error) {
	// check that hostname matches policy hostname.namespace.k8s.zach
	regex := fmt.Sprintf(`.*[\w-]+\.%s.%s`, regexp.QuoteMeta(namespace), regexp.QuoteMeta(baseDomain))
	match, _ := regexp.MatchString(regex, hostname)

	// use external dns client to check for duplicate A records

	return match, nil
}
