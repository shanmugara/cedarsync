package controller

import (
	"context"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
	"time"

	cedarsyncv1alpha1 "github.com/shanmugara/cedarsync/api/v1alpha1"
)

// CedarPolicyReconciler reconciles a CedarPolicy object

const (
	PolicyFinalizer = "cedarpolicy.cedarsync.omegahome.net/finalizer"
	cedarAPIUrl     = "http://127.0.0.1:5000"
)

type PolicyBlock struct {
	Id          int    `json:"id"`
	Principal   string `json:"principal"`
	Action      string `json:"action"`
	Resource    string `json:"resource"`
	Effect      string `json:"effect"`
	Condition   string `json:"condition"`
	Annotations string `json:"annotations"`
}

// ReconcilePolicy is the function to reconcile the CedarPolicy object

func (r *CedarPolicyReconciler) ReconcilePolicy(ctx context.Context, capi *cedarsyncv1alpha1.CedarApi) error {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling CedarPolicy")

	//Fetch the CedarPolicy instance for the CedarApi instance
	cedarPolicy := &cedarsyncv1alpha1.CedarPolicy{}
	err := r.Get(ctx, client.ObjectKey{Namespace: capi.Namespace, Name: capi.Spec.Cluster}, cedarPolicy)

	if errors.IsNotFound(err) {
		logger.Info("CedarPolicy not found. Creating a new one")
		//Create a new CedarPolicy instance

		policy, err := r.FetchCedarPolicy(ctx, capi)
		if err != nil {
			logger.Error(err, "Failed to fetch CedarPolicy")
			return err
		}

		cedarPolicy = &cedarsyncv1alpha1.CedarPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      capi.Spec.Cluster,
				Namespace: capi.Namespace,
			},
			Spec: cedarsyncv1alpha1.CedarPolicySpec{
				Policy: policy,
			},
		}
		if err := controllerutil.SetOwnerReference(capi, cedarPolicy, r.Scheme); err != nil {
			logger.Error(err, "Failed to set owner reference")
			return err
		}

		if err := r.Create(ctx, cedarPolicy); err != nil {
			logger.Error(err, "Failed to create CedarPolicy")
			return err
		}

		//CedarPolicy created successfully
		logger.Info("CedarPolicy created successfully")
		return nil

	}

	return nil

}

func (r *CedarPolicyReconciler) DeletePolicy(ctx context.Context, capi *cedarsyncv1alpha1.CedarApi) error {
	logger := log.FromContext(ctx)

	logger.Info("Deleting CedarPolicy")

	//Fetch the CedarPolicy instance for the CedarApi instance
	cedarPolicy := &cedarsyncv1alpha1.CedarPolicy{}
	if err := r.Get(ctx, client.ObjectKey{Namespace: capi.Namespace, Name: capi.Spec.Cluster}, cedarPolicy); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("CedarPolicy not found. Ignoring since object must be deleted")
			return nil
		}
		logger.Error(err, "Failed to fetch CedarPolicy")
		return err
	}

	//delete the CedarPolicy instance
	if err := r.Delete(ctx, cedarPolicy); err != nil {
		logger.Error(err, "Failed to delete CedarPolicy")
		return err
	}

	return nil
}

func (r *CedarPolicyReconciler) FetchCedarPolicy(ctx context.Context, capi *cedarsyncv1alpha1.CedarApi) (*cedarsyncv1alpha1.Policy, error) {
	logger := log.FromContext(ctx)

	logger.Info("Fetching CedarPolicy from cedar API")

	client := &http.Client{Timeout: time.Second * 2}
	callURI := capi.Spec.ApiUrl + "/api/cluster/" + capi.Spec.Cluster
	req, err := http.NewRequest("GET", callURI, nil)
	if err != nil {
		logger.Error(err, "Failed to create new request")
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err, "Failed to fetch CedarPolicy")
		return nil, err
	}
	defer resp.Body.Close()

	//Check Status Code

	switch resp.StatusCode {
	case http.StatusNotFound:
		logger.Info("CedarPolicy not found for the cluster")
		return nil, err
	case http.StatusOK:
		logger.Info("CedarPolicy found for the cluster")
		break
	default:
		logger.Error(err, "Failed to fetch CedarPolicy")
		return nil, err
	}

	//Parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err, "Failed to read response body")
		return nil, err
	}

	var apiResponse PolicyBlock
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		logger.Error(err, "Failed to unmarshal response body")
		return nil, err
	}

	var annotations map[string]string
	if err := json.Unmarshal([]byte(apiResponse.Annotations), &annotations); err != nil {
		logger.Error(err, "Failed to unmarshal annotations")
		return nil, err
	}

	var policy cedarsyncv1alpha1.Policy
	policy.Principal = apiResponse.Principal
	policy.Action = apiResponse.Action
	policy.Resource = apiResponse.Resource
	policy.Effect = apiResponse.Effect
	policy.Conditions = strings.Split(apiResponse.Condition, ",")
	policy.Annotations = annotations

	//Fetch the CedarPolicy instance for the CedarApi instance
	return &policy, nil
}
