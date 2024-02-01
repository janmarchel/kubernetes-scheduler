package pkg

import (
	"encoding/json"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/api/core/v1"
	"net/http"
)

func HandleMutate(w http.ResponseWriter, r *http.Request) {
	// Odczytaj ciało żądania
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}

	// Parsowanie AdmissionReview żądania
	var admissionReviewReq v1beta1.AdmissionReview
	if err := json.Unmarshal(body, &admissionReviewReq); err != nil {
		http.Error(w, "could not unmarshal request", http.StatusBadRequest)
		return
	}

	// Przygotowanie odpowiedzi
	admissionReviewResponse := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			UID: admissionReviewReq.Request.UID,
		},
	}

	// Deserializacja Pod z AdmissionRequest
	var pod v1.Pod
	if err := json.Unmarshal(admissionReviewReq.Request.Object.Raw, &pod); err != nil {
		admissionReviewResponse.Response.Allowed = false
		admissionReviewResponse.Response.Result.Message = "could not deserialize pod"
	} else {
		// Tutaj dokonujemy mutacji poda
		// Przykład: zmiana etykiety
		if pod.Labels == nil {
			pod.Labels = map[string]string{}
		}
		pod.Labels["custom-label"] = "custom-value"

		// Serializacja zmodyfikowanego Pod do JSON
		modifiedPodJSON, err := json.Marshal(pod)
		if err != nil {
			admissionReviewResponse.Response.Allowed = false
			admissionReviewResponse.Response.Result.Message = "could not serialize modified pod"
		} else {
			// Stwórz JSONPatch
			patch := []map[string]string{{
				"op":    "add",
				"path":  "/metadata/labels/custom-label",
				"value": "custom-value",
			}}
			patchBytes, _ := json.Marshal(patch)

			admissionReviewResponse.Response.Allowed = true
			admissionReviewResponse.Response.Patch = patchBytes
			admissionReviewResponse.Response.PatchType = new(v1beta1.PatchType)
			*admissionReviewResponse.Response.PatchType = v1beta1.PatchTypeJSONPatch
		}
	}

	// Serializacja odpowiedzi AdmissionReview do JSON i wysłanie
	respBytes, err := json.Marshal(admissionReviewResponse)
	if err != nil {
		http.Error(w, "could not serialize response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respBytes); err != nil {
		http.Error(w, "could not write response", http.StatusInternalServerError)
		return
	}
}
