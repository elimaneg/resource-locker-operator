package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// ResourceLockerSpec defines the desired state of ResourceLocker
// +k8s:openapi-gen=true
type ResourceLockerSpec struct {

	// Resources is a list of resource manifests that should be locked into the specified configuration
	// +kubebuilder:validation:Optional
	// +listType=atomic
	Resources []Resource `json:"resources,omitempty"`

	// Patches is a list of pacthes that should be encforced at runtime.
	// +kubebuilder:validation:Optional
	// +listType=atomic
	Patches []Patch `json:"patches,omitempty"`

	// ServiceAccountRef is the service account to be used to run the controllers associated with this configuration
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="{Name:default}"
	ServiceAccountRef corev1.LocalObjectReference `json:"serviceAccountRef,omitempty"`
}

// +k8s:openapi-gen=true
type Resource struct {
	// +kubebuilder:validation:Required
	Object runtime.RawExtension `json:"object"`
	// +kubebuilder:validation:Optional
	// +listType=set
	ExcludedPaths []string `json:"excludedPaths,omitempty"`
}

// Patch describe a pacth to be enforced at runtim
// +k8s:openapi-gen=true
type Patch struct {
	// SourceObject refs is an arrays of refereces to source objects that will be used as input for the template processing
	// +kubebuilder:validation:Optional
	// +listType=atomic
	SourceObjectRefs []corev1.ObjectReference `json:"sourceObjectRefs,omitempty"`

	// TargetObjectRef is a reference to the object to which the pacth should be applied.
	// +kubebuilder:validation:Required
	TargetObjectRef corev1.ObjectReference `json:"targetObjectRef"`

	// PatchType is the type of patch to be applied, one of "application/json-patch+json"'"application/merge-patch+json","application/strategic-merge-patch+json","application/apply-patch+yaml"
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum="application/json-patch+json";"application/merge-patch+json";"application/strategic-merge-patch+json";"application/apply-patch+yaml"
	// +kubebuilder:default:="application/strategic-merge-patch+json"
	PatchType types.PatchType `json:"patchType,omitempty"`

	// PatchTemplate is a go template that will be resolved using the SourceObjectRefs as parameters. The result must be a valid patch based on the pacth type and the target object.
	// +kubebuilder:validation:Required
	PatchTemplate string `json:"patchTemplate"`
}

// ResourceLockerStatus defines the observed state of ResourceLocker
// +k8s:openapi-gen=true
type ResourceLockerStatus struct {
	// Type of deployment condition.
	// +kubebuilder:validation:Enum=Success;Failure
	// +kubebuilder:validation:Required
	Type ConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=True;False;Unknown
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	// +kubebuilder:validation:Optional
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// A human readable message indicating details about the transition.
	// +kubebuilder:validation:Optional
	Message string `json:"message,omitempty"`
	// +kubebuilder:validation:Optional
	// +listType=atomic
	ResourceStatuses []LockingStatus `json:"resourceStatuses,omitempty"`
	// +kubebuilder:validation:Optional
	// +listType=atomic
	PatchStatuses []LockingStatus `json:"patchStatuses,omitempty"`
}

// +k8s:openapi-gen=true
type LockingStatus struct {
	// the name of the locked configuration
	// +kubebuilder:validation:Required
	ObjectReference corev1.ObjectReference `json:"objectReference"`
	// Type of deployment condition.
	// +kubebuilder:validation:Enum=Enforcing;Failure
	// +kubebuilder:validation:Optional
	Type ConditionType `json:"type,omitempty"`
	// Status of the condition, one of True, False, Unknown.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=True;False;Unknown
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	// +kubebuilder:validation:Optional
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// A human readable message indicating details about the transition.
	// +kubebuilder:validation:Optional
	Message string `json:"message,omitempty"`
}

type ConditionType string

const (
	// Enforcing means that the resource or patch has been succesfully reconciled and it's being enforced
	Enforcing ConditionType = "Enforcing"
	// Failure means that the resource or patch has not been successfully reconciled and we cannot guarntee that it's being enforced
	Failure ConditionType = "Failure"
	// Success means that the instance has been successfully reconciled
	Success ConditionType = "Success"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourceLocker is the Schema for the resourcelockers API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=resourcelockers,scope=Namespaced
type ResourceLocker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceLockerSpec   `json:"spec,omitempty"`
	Status ResourceLockerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourceLockerList contains a list of ResourceLocker
type ResourceLockerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceLocker `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceLocker{}, &ResourceLockerList{})
}
