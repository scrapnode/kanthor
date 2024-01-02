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

// checks if the EndpointUpdateReq type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &EndpointUpdateReq{}

// EndpointUpdateReq struct for EndpointUpdateReq
type EndpointUpdateReq struct {
	Method *string `json:"method,omitempty"`
	Name *string `json:"name,omitempty"`
}

// NewEndpointUpdateReq instantiates a new EndpointUpdateReq object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEndpointUpdateReq() *EndpointUpdateReq {
	this := EndpointUpdateReq{}
	return &this
}

// NewEndpointUpdateReqWithDefaults instantiates a new EndpointUpdateReq object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEndpointUpdateReqWithDefaults() *EndpointUpdateReq {
	this := EndpointUpdateReq{}
	return &this
}

// GetMethod returns the Method field value if set, zero value otherwise.
func (o *EndpointUpdateReq) GetMethod() string {
	if o == nil || IsNil(o.Method) {
		var ret string
		return ret
	}
	return *o.Method
}

// GetMethodOk returns a tuple with the Method field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointUpdateReq) GetMethodOk() (*string, bool) {
	if o == nil || IsNil(o.Method) {
		return nil, false
	}
	return o.Method, true
}

// HasMethod returns a boolean if a field has been set.
func (o *EndpointUpdateReq) HasMethod() bool {
	if o != nil && !IsNil(o.Method) {
		return true
	}

	return false
}

// SetMethod gets a reference to the given string and assigns it to the Method field.
func (o *EndpointUpdateReq) SetMethod(v string) {
	o.Method = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *EndpointUpdateReq) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointUpdateReq) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *EndpointUpdateReq) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *EndpointUpdateReq) SetName(v string) {
	o.Name = &v
}

func (o EndpointUpdateReq) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o EndpointUpdateReq) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Method) {
		toSerialize["method"] = o.Method
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	return toSerialize, nil
}

type NullableEndpointUpdateReq struct {
	value *EndpointUpdateReq
	isSet bool
}

func (v NullableEndpointUpdateReq) Get() *EndpointUpdateReq {
	return v.value
}

func (v *NullableEndpointUpdateReq) Set(val *EndpointUpdateReq) {
	v.value = val
	v.isSet = true
}

func (v NullableEndpointUpdateReq) IsSet() bool {
	return v.isSet
}

func (v *NullableEndpointUpdateReq) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEndpointUpdateReq(val *EndpointUpdateReq) *NullableEndpointUpdateReq {
	return &NullableEndpointUpdateReq{value: val, isSet: true}
}

func (v NullableEndpointUpdateReq) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEndpointUpdateReq) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


