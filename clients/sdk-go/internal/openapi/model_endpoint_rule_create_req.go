/*
Kanthor SDK API

SDK API

API version: 1.0
Contact: support@kanthorlabs.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the EndpointRuleCreateReq type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &EndpointRuleCreateReq{}

// EndpointRuleCreateReq struct for EndpointRuleCreateReq
type EndpointRuleCreateReq struct {
	ConditionExpression *string `json:"condition_expression,omitempty"`
	ConditionSource *string `json:"condition_source,omitempty"`
	EpId *string `json:"ep_id,omitempty"`
	Exclusionary *bool `json:"exclusionary,omitempty"`
	Name *string `json:"name,omitempty"`
	Priority *int32 `json:"priority,omitempty"`
}

// NewEndpointRuleCreateReq instantiates a new EndpointRuleCreateReq object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEndpointRuleCreateReq() *EndpointRuleCreateReq {
	this := EndpointRuleCreateReq{}
	return &this
}

// NewEndpointRuleCreateReqWithDefaults instantiates a new EndpointRuleCreateReq object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEndpointRuleCreateReqWithDefaults() *EndpointRuleCreateReq {
	this := EndpointRuleCreateReq{}
	return &this
}

// GetConditionExpression returns the ConditionExpression field value if set, zero value otherwise.
func (o *EndpointRuleCreateReq) GetConditionExpression() string {
	if o == nil || IsNil(o.ConditionExpression) {
		var ret string
		return ret
	}
	return *o.ConditionExpression
}

// GetConditionExpressionOk returns a tuple with the ConditionExpression field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointRuleCreateReq) GetConditionExpressionOk() (*string, bool) {
	if o == nil || IsNil(o.ConditionExpression) {
		return nil, false
	}
	return o.ConditionExpression, true
}

// HasConditionExpression returns a boolean if a field has been set.
func (o *EndpointRuleCreateReq) HasConditionExpression() bool {
	if o != nil && !IsNil(o.ConditionExpression) {
		return true
	}

	return false
}

// SetConditionExpression gets a reference to the given string and assigns it to the ConditionExpression field.
func (o *EndpointRuleCreateReq) SetConditionExpression(v string) {
	o.ConditionExpression = &v
}

// GetConditionSource returns the ConditionSource field value if set, zero value otherwise.
func (o *EndpointRuleCreateReq) GetConditionSource() string {
	if o == nil || IsNil(o.ConditionSource) {
		var ret string
		return ret
	}
	return *o.ConditionSource
}

// GetConditionSourceOk returns a tuple with the ConditionSource field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointRuleCreateReq) GetConditionSourceOk() (*string, bool) {
	if o == nil || IsNil(o.ConditionSource) {
		return nil, false
	}
	return o.ConditionSource, true
}

// HasConditionSource returns a boolean if a field has been set.
func (o *EndpointRuleCreateReq) HasConditionSource() bool {
	if o != nil && !IsNil(o.ConditionSource) {
		return true
	}

	return false
}

// SetConditionSource gets a reference to the given string and assigns it to the ConditionSource field.
func (o *EndpointRuleCreateReq) SetConditionSource(v string) {
	o.ConditionSource = &v
}

// GetEpId returns the EpId field value if set, zero value otherwise.
func (o *EndpointRuleCreateReq) GetEpId() string {
	if o == nil || IsNil(o.EpId) {
		var ret string
		return ret
	}
	return *o.EpId
}

// GetEpIdOk returns a tuple with the EpId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointRuleCreateReq) GetEpIdOk() (*string, bool) {
	if o == nil || IsNil(o.EpId) {
		return nil, false
	}
	return o.EpId, true
}

// HasEpId returns a boolean if a field has been set.
func (o *EndpointRuleCreateReq) HasEpId() bool {
	if o != nil && !IsNil(o.EpId) {
		return true
	}

	return false
}

// SetEpId gets a reference to the given string and assigns it to the EpId field.
func (o *EndpointRuleCreateReq) SetEpId(v string) {
	o.EpId = &v
}

// GetExclusionary returns the Exclusionary field value if set, zero value otherwise.
func (o *EndpointRuleCreateReq) GetExclusionary() bool {
	if o == nil || IsNil(o.Exclusionary) {
		var ret bool
		return ret
	}
	return *o.Exclusionary
}

// GetExclusionaryOk returns a tuple with the Exclusionary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointRuleCreateReq) GetExclusionaryOk() (*bool, bool) {
	if o == nil || IsNil(o.Exclusionary) {
		return nil, false
	}
	return o.Exclusionary, true
}

// HasExclusionary returns a boolean if a field has been set.
func (o *EndpointRuleCreateReq) HasExclusionary() bool {
	if o != nil && !IsNil(o.Exclusionary) {
		return true
	}

	return false
}

// SetExclusionary gets a reference to the given bool and assigns it to the Exclusionary field.
func (o *EndpointRuleCreateReq) SetExclusionary(v bool) {
	o.Exclusionary = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *EndpointRuleCreateReq) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointRuleCreateReq) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *EndpointRuleCreateReq) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *EndpointRuleCreateReq) SetName(v string) {
	o.Name = &v
}

// GetPriority returns the Priority field value if set, zero value otherwise.
func (o *EndpointRuleCreateReq) GetPriority() int32 {
	if o == nil || IsNil(o.Priority) {
		var ret int32
		return ret
	}
	return *o.Priority
}

// GetPriorityOk returns a tuple with the Priority field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointRuleCreateReq) GetPriorityOk() (*int32, bool) {
	if o == nil || IsNil(o.Priority) {
		return nil, false
	}
	return o.Priority, true
}

// HasPriority returns a boolean if a field has been set.
func (o *EndpointRuleCreateReq) HasPriority() bool {
	if o != nil && !IsNil(o.Priority) {
		return true
	}

	return false
}

// SetPriority gets a reference to the given int32 and assigns it to the Priority field.
func (o *EndpointRuleCreateReq) SetPriority(v int32) {
	o.Priority = &v
}

func (o EndpointRuleCreateReq) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o EndpointRuleCreateReq) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ConditionExpression) {
		toSerialize["condition_expression"] = o.ConditionExpression
	}
	if !IsNil(o.ConditionSource) {
		toSerialize["condition_source"] = o.ConditionSource
	}
	if !IsNil(o.EpId) {
		toSerialize["ep_id"] = o.EpId
	}
	if !IsNil(o.Exclusionary) {
		toSerialize["exclusionary"] = o.Exclusionary
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Priority) {
		toSerialize["priority"] = o.Priority
	}
	return toSerialize, nil
}

type NullableEndpointRuleCreateReq struct {
	value *EndpointRuleCreateReq
	isSet bool
}

func (v NullableEndpointRuleCreateReq) Get() *EndpointRuleCreateReq {
	return v.value
}

func (v *NullableEndpointRuleCreateReq) Set(val *EndpointRuleCreateReq) {
	v.value = val
	v.isSet = true
}

func (v NullableEndpointRuleCreateReq) IsSet() bool {
	return v.isSet
}

func (v *NullableEndpointRuleCreateReq) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEndpointRuleCreateReq(val *EndpointRuleCreateReq) *NullableEndpointRuleCreateReq {
	return &NullableEndpointRuleCreateReq{value: val, isSet: true}
}

func (v NullableEndpointRuleCreateReq) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEndpointRuleCreateReq) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

