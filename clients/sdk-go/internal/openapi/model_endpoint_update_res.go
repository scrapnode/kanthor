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

// checks if the EndpointUpdateRes type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &EndpointUpdateRes{}

// EndpointUpdateRes struct for EndpointUpdateRes
type EndpointUpdateRes struct {
	AppId *string `json:"app_id,omitempty"`
	CreatedAt *int32 `json:"created_at,omitempty"`
	Id *string `json:"id,omitempty"`
	Method *string `json:"method,omitempty"`
	Name *string `json:"name,omitempty"`
	UpdatedAt *int32 `json:"updated_at,omitempty"`
	Uri *string `json:"uri,omitempty"`
}

// NewEndpointUpdateRes instantiates a new EndpointUpdateRes object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEndpointUpdateRes() *EndpointUpdateRes {
	this := EndpointUpdateRes{}
	return &this
}

// NewEndpointUpdateResWithDefaults instantiates a new EndpointUpdateRes object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEndpointUpdateResWithDefaults() *EndpointUpdateRes {
	this := EndpointUpdateRes{}
	return &this
}

// GetAppId returns the AppId field value if set, zero value otherwise.
func (o *EndpointUpdateRes) GetAppId() string {
	if o == nil || IsNil(o.AppId) {
		var ret string
		return ret
	}
	return *o.AppId
}

// GetAppIdOk returns a tuple with the AppId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointUpdateRes) GetAppIdOk() (*string, bool) {
	if o == nil || IsNil(o.AppId) {
		return nil, false
	}
	return o.AppId, true
}

// HasAppId returns a boolean if a field has been set.
func (o *EndpointUpdateRes) HasAppId() bool {
	if o != nil && !IsNil(o.AppId) {
		return true
	}

	return false
}

// SetAppId gets a reference to the given string and assigns it to the AppId field.
func (o *EndpointUpdateRes) SetAppId(v string) {
	o.AppId = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise.
func (o *EndpointUpdateRes) GetCreatedAt() int32 {
	if o == nil || IsNil(o.CreatedAt) {
		var ret int32
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointUpdateRes) GetCreatedAtOk() (*int32, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}
	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *EndpointUpdateRes) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given int32 and assigns it to the CreatedAt field.
func (o *EndpointUpdateRes) SetCreatedAt(v int32) {
	o.CreatedAt = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *EndpointUpdateRes) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointUpdateRes) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *EndpointUpdateRes) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *EndpointUpdateRes) SetId(v string) {
	o.Id = &v
}

// GetMethod returns the Method field value if set, zero value otherwise.
func (o *EndpointUpdateRes) GetMethod() string {
	if o == nil || IsNil(o.Method) {
		var ret string
		return ret
	}
	return *o.Method
}

// GetMethodOk returns a tuple with the Method field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointUpdateRes) GetMethodOk() (*string, bool) {
	if o == nil || IsNil(o.Method) {
		return nil, false
	}
	return o.Method, true
}

// HasMethod returns a boolean if a field has been set.
func (o *EndpointUpdateRes) HasMethod() bool {
	if o != nil && !IsNil(o.Method) {
		return true
	}

	return false
}

// SetMethod gets a reference to the given string and assigns it to the Method field.
func (o *EndpointUpdateRes) SetMethod(v string) {
	o.Method = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *EndpointUpdateRes) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointUpdateRes) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *EndpointUpdateRes) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *EndpointUpdateRes) SetName(v string) {
	o.Name = &v
}

// GetUpdatedAt returns the UpdatedAt field value if set, zero value otherwise.
func (o *EndpointUpdateRes) GetUpdatedAt() int32 {
	if o == nil || IsNil(o.UpdatedAt) {
		var ret int32
		return ret
	}
	return *o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointUpdateRes) GetUpdatedAtOk() (*int32, bool) {
	if o == nil || IsNil(o.UpdatedAt) {
		return nil, false
	}
	return o.UpdatedAt, true
}

// HasUpdatedAt returns a boolean if a field has been set.
func (o *EndpointUpdateRes) HasUpdatedAt() bool {
	if o != nil && !IsNil(o.UpdatedAt) {
		return true
	}

	return false
}

// SetUpdatedAt gets a reference to the given int32 and assigns it to the UpdatedAt field.
func (o *EndpointUpdateRes) SetUpdatedAt(v int32) {
	o.UpdatedAt = &v
}

// GetUri returns the Uri field value if set, zero value otherwise.
func (o *EndpointUpdateRes) GetUri() string {
	if o == nil || IsNil(o.Uri) {
		var ret string
		return ret
	}
	return *o.Uri
}

// GetUriOk returns a tuple with the Uri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointUpdateRes) GetUriOk() (*string, bool) {
	if o == nil || IsNil(o.Uri) {
		return nil, false
	}
	return o.Uri, true
}

// HasUri returns a boolean if a field has been set.
func (o *EndpointUpdateRes) HasUri() bool {
	if o != nil && !IsNil(o.Uri) {
		return true
	}

	return false
}

// SetUri gets a reference to the given string and assigns it to the Uri field.
func (o *EndpointUpdateRes) SetUri(v string) {
	o.Uri = &v
}

func (o EndpointUpdateRes) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o EndpointUpdateRes) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.AppId) {
		toSerialize["app_id"] = o.AppId
	}
	if !IsNil(o.CreatedAt) {
		toSerialize["created_at"] = o.CreatedAt
	}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Method) {
		toSerialize["method"] = o.Method
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.UpdatedAt) {
		toSerialize["updated_at"] = o.UpdatedAt
	}
	if !IsNil(o.Uri) {
		toSerialize["uri"] = o.Uri
	}
	return toSerialize, nil
}

type NullableEndpointUpdateRes struct {
	value *EndpointUpdateRes
	isSet bool
}

func (v NullableEndpointUpdateRes) Get() *EndpointUpdateRes {
	return v.value
}

func (v *NullableEndpointUpdateRes) Set(val *EndpointUpdateRes) {
	v.value = val
	v.isSet = true
}

func (v NullableEndpointUpdateRes) IsSet() bool {
	return v.isSet
}

func (v *NullableEndpointUpdateRes) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEndpointUpdateRes(val *EndpointUpdateRes) *NullableEndpointUpdateRes {
	return &NullableEndpointUpdateRes{value: val, isSet: true}
}

func (v NullableEndpointUpdateRes) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEndpointUpdateRes) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


