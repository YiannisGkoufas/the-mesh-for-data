// Copyright 2020 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CatalogRequirements contain the specifics for catalogging the data asset
type CatalogRequirements struct {
	// CatalogService specifies the datacatalog service that will be used for catalogging the data into.
	// +optional
	CatalogService string `json:"service,omitempty"`

	// CatalogID specifies the catalog where the data will be cataloged.
	// +optional
	CatalogID string `json:"catalogID,omitempty"`
}

// CopyRequirements include the requirements for the data copy operation
type CopyRequirements struct {
	// Required indicates that the data must be copied.
	// +optional
	Required bool `json:"required,omitempty"`

	// Catalog indicates that the data asset must be cataloged.
	// +optional
	Catalog CatalogRequirements `json:"catalog,omitempty"`
}

// DataRequirements structure contains a list of requirements (interface, need to catalog the dataset, etc.)
type DataRequirements struct {
	// Interface indicates the protocol and format expected by the data user
	// +required
	Interface InterfaceDetails `json:"interface"`

	// CopyRequrements include the requirements for copying the data
	// +optional
	Copy CopyRequirements `json:"copy,omitempty"`
}

// DataContext indicates data set chosen by the Data Scientist to be used by his application,
// and includes information about the data format and technologies used by the application
// to access the data.
type DataContext struct {
	// DataSetID is a unique identifier of the dataset chosen from the data catalog for processing by the data user application.
	// +required
	// +kubebuilder:validation:MinLength=1
	DataSetID string `json:"dataSetID"`

	// CatalogService represents the catalog service for accessing the requested dataset.
	// If not specified, the enterprise catalog service will be used.
	// +optional
	CatalogService string `json:"catalogService,omitempty"`

	// Requirements from the system
	// +required
	Requirements DataRequirements `json:"requirements"`
}

// AppUserRole indicates the role required to use the application
type AppUserRole string

// ApplicationDetails provides information about the Data Scientist's application, which is deployed separately.
// The information provided is used to determine if the data should be altered in any way prior to its use,
// based on policies and rules defined in an external data policy manager.
type ApplicationDetails struct {
	// Purpose indicates the reason for the processing and the use of the data by the Data Scientist's application.
	// +required
	Purpose string `json:"purpose,omitempty"`

	// Role indicates the position held or role filled by the Data Scientist as it relates to the processing of the
	// data he has chosen.
	// +required
	Role AppUserRole `json:"role"`
}

// M4DApplicationSpec defines the desired state of M4DApplication.
type M4DApplicationSpec struct {

	// Selector enables to connect the resource to the application
	// Application labels should match the labels in the selector.
	// For some flows the selector may not be used.
	// +optional
	Selector Selector `json:"selector"`

	// AppInfo contains information describing the reasons and geography of the processing
	// that will be done by the Data Scientist's application.
	// +required
	AppInfo ApplicationDetails `json:"appInfo"`

	// Data contains the identifiers of the data to be used by the Data Scientist's application,
	// and the protocol used to access it and the format expected.
	// +required
	// +kubebuilder:validation:MinItems=1
	Data []DataContext `json:"data"`
}

// ErrorMessages that are reported to the user
const (
	ReadAccessDenied            string = "Governance policies forbid access to the data."
	CopyNotAllowed              string = "Copy of the data is required but can not be done according to the governance policies."
	WriteNotAllowed             string = "Governance policies forbid writing of the data."
	ModuleNotFound              string = "No module has been registered"
	InsufficientStorage         string = "No bucket was provisioned for implicit copy"
	InvalidClusterConfiguration string = "Cluster configuration does not support the requirements."
)

// Condition indices are static. Conditions always present in the status.
const (
	FailureConditionIndex int64 = 0
	ErrorConditionIndex   int64 = 1
)

// ConditionType represents a condition type
type ConditionType string

const (
	// ErrorCondition means that an error was encountered during blueprint construction
	ErrorCondition ConditionType = "Error"

	// FailureCondition means that a blueprint could not be constructed
	FailureCondition ConditionType = "Failure"
)

// Condition describes the state of a M4DApplication at a certain point.
type Condition struct {
	// Type of the condition
	Type ConditionType `json:"type"`
	// Status of the condition: true or false
	Status corev1.ConditionStatus `json:"status"`
	// Message contains the details of the current condition
	// +optional
	Message string `json:"message,omitempty"`
}

// ResourceReference contains resource identifier(name, namespace, kind)
type ResourceReference struct {
	// Name of the resource
	Name string `json:"name"`
	// Namespace of the resource
	Namespace string `json:"namespace"`
	// Kind of the resource (Blueprint, Plotter)
	Kind string `json:"kind"`
}

// M4DApplicationStatus defines the observed state of M4DApplication.
type M4DApplicationStatus struct {

	// Ready is true if a blueprint has been successfully orchestrated
	Ready bool `json:"ready,omitempty"`

	// Conditions represent the possible error and failure conditions
	// +optional
	Conditions []Condition `json:"conditions,omitempty"`

	// DataAccessInstructions indicate how the data user or his application may access the data.
	// Instructions are available upon successful orchestration.
	// +optional
	DataAccessInstructions string `json:"dataAccessInstructions,omitempty"`

	// ObservedGeneration is taken from the M4DApplication metadata.  This is used to determine during reconcile
	// whether reconcile was called because the desired state changed, or whether the Blueprint status changed.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Generated resource identifier
	// +optional
	Generated *ResourceReference `json:"generated,omitempty"`
}

// M4DApplication provides information about the application being used by a Data Scientist,
// the nature of the processing, and the data sets that the Data Scientist has chosen for processing by the application.
// The M4DApplication controller (aka pilot) obtains instructions regarding any governance related changes that must
// be performed on the data, identifies the modules capable of performing such changes, and finally
// generates the Blueprint which defines the secure runtime environment and all the components
// in it.  This runtime environment provides the Data Scientist's application with access to the data requested
// in a secure manner and without having to provide any credentials for the data sets.  The credentials are obtained automatically
// by the manager from an external credential management system, which may or may not be part of a data catalog.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type M4DApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   M4DApplicationSpec   `json:"spec,omitempty"`
	Status M4DApplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// M4DApplicationList contains a list of M4DApplication
type M4DApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []M4DApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&M4DApplication{}, &M4DApplicationList{})
}
