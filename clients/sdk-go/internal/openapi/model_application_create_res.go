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

// checks if the ApplicationCreateRes type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ApplicationCreateRes{}

// ApplicationCreateRes struct for ApplicationCreateRes
type ApplicationCreateRes struct {
	CreatedAt *int32 `json:"created_at,omitempty"`
	Id *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	UpdatedAt *int32 `json:"updated_at,omitempty"`
	WsId *string `json:"ws_id,omitempty"`
}

// NewApplicationCreateRes instantiates a new ApplicationCreateRes object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApplicationCreateRes() *ApplicationCreateRes {
	this := ApplicationCreateRes{}
	return &this
}

// NewApplicationCreateResWithDefaults instantiates a new ApplicationCreateRes object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApplicationCreateResWithDefaults() *ApplicationCreateRes {
	this := ApplicationCreateRes{}
	return &this
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise.
func (o *ApplicationCreateRes) GetCreatedAt() int32 {
	if o == nil || IsNil(o.CreatedAt) {
		var ret int32
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApplicationCreateRes) GetCreatedAtOk() (*int32, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}
	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *ApplicationCreateRes) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given int32 and assigns it to the CreatedAt field.
func (o *ApplicationCreateRes) SetCreatedAt(v int32) {
	o.CreatedAt = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *ApplicationCreateRes) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApplicationCreateRes) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ApplicationCreateRes) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ApplicationCreateRes) SetId(v string) {
	o.Id = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ApplicationCreateRes) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApplicationCreateRes) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ApplicationCreateRes) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ApplicationCreateRes) SetName(v string) {
	o.Name = &v
}

// GetUpdatedAt returns the UpdatedAt field value if set, zero value otherwise.
func (o *ApplicationCreateRes) GetUpdatedAt() int32 {
	if o == nil || IsNil(o.UpdatedAt) {
		var ret int32
		return ret
	}
	return *o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApplicationCreateRes) GetUpdatedAtOk() (*int32, bool) {
	if o == nil || IsNil(o.UpdatedAt) {
		return nil, false
	}
	return o.UpdatedAt, true
}

// HasUpdatedAt returns a boolean if a field has been set.
func (o *ApplicationCreateRes) HasUpdatedAt() bool {
	if o != nil && !IsNil(o.UpdatedAt) {
		return true
	}

	return false
}

// SetUpdatedAt gets a reference to the given int32 and assigns it to the UpdatedAt field.
func (o *ApplicationCreateRes) SetUpdatedAt(v int32) {
	o.UpdatedAt = &v
}

// GetWsId returns the WsId field value if set, zero value otherwise.
func (o *ApplicationCreateRes) GetWsId() string {
	if o == nil || IsNil(o.WsId) {
		var ret string
		return ret
	}
	return *o.WsId
}

// GetWsIdOk returns a tuple with the WsId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApplicationCreateRes) GetWsIdOk() (*string, bool) {
	if o == nil || IsNil(o.WsId) {
		return nil, false
	}
	return o.WsId, true
}

// HasWsId returns a boolean if a field has been set.
func (o *ApplicationCreateRes) HasWsId() bool {
	if o != nil && !IsNil(o.WsId) {
		return true
	}

	return false
}

// SetWsId gets a reference to the given string and assigns it to the WsId field.
func (o *ApplicationCreateRes) SetWsId(v string) {
	o.WsId = &v
}

func (o ApplicationCreateRes) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ApplicationCreateRes) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.CreatedAt) {
		toSerialize["created_at"] = o.CreatedAt
	}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.UpdatedAt) {
		toSerialize["updated_at"] = o.UpdatedAt
	}
	if !IsNil(o.WsId) {
		toSerialize["ws_id"] = o.WsId
	}
	return toSerialize, nil
}

type NullableApplicationCreateRes struct {
	value *ApplicationCreateRes
	isSet bool
}

func (v NullableApplicationCreateRes) Get() *ApplicationCreateRes {
	return v.value
}

func (v *NullableApplicationCreateRes) Set(val *ApplicationCreateRes) {
	v.value = val
	v.isSet = true
}

func (v NullableApplicationCreateRes) IsSet() bool {
	return v.isSet
}

func (v *NullableApplicationCreateRes) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApplicationCreateRes(val *ApplicationCreateRes) *NullableApplicationCreateRes {
	return &NullableApplicationCreateRes{value: val, isSet: true}
}

func (v NullableApplicationCreateRes) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApplicationCreateRes) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


