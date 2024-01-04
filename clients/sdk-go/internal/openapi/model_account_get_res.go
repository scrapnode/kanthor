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

// checks if the AccountGetRes type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AccountGetRes{}

// AccountGetRes struct for AccountGetRes
type AccountGetRes struct {
	Account *AuthenticatorAccount `json:"account,omitempty"`
}

// NewAccountGetRes instantiates a new AccountGetRes object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAccountGetRes() *AccountGetRes {
	this := AccountGetRes{}
	return &this
}

// NewAccountGetResWithDefaults instantiates a new AccountGetRes object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAccountGetResWithDefaults() *AccountGetRes {
	this := AccountGetRes{}
	return &this
}

// GetAccount returns the Account field value if set, zero value otherwise.
func (o *AccountGetRes) GetAccount() AuthenticatorAccount {
	if o == nil || IsNil(o.Account) {
		var ret AuthenticatorAccount
		return ret
	}
	return *o.Account
}

// GetAccountOk returns a tuple with the Account field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AccountGetRes) GetAccountOk() (*AuthenticatorAccount, bool) {
	if o == nil || IsNil(o.Account) {
		return nil, false
	}
	return o.Account, true
}

// HasAccount returns a boolean if a field has been set.
func (o *AccountGetRes) HasAccount() bool {
	if o != nil && !IsNil(o.Account) {
		return true
	}

	return false
}

// SetAccount gets a reference to the given AuthenticatorAccount and assigns it to the Account field.
func (o *AccountGetRes) SetAccount(v AuthenticatorAccount) {
	o.Account = &v
}

func (o AccountGetRes) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AccountGetRes) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Account) {
		toSerialize["account"] = o.Account
	}
	return toSerialize, nil
}

type NullableAccountGetRes struct {
	value *AccountGetRes
	isSet bool
}

func (v NullableAccountGetRes) Get() *AccountGetRes {
	return v.value
}

func (v *NullableAccountGetRes) Set(val *AccountGetRes) {
	v.value = val
	v.isSet = true
}

func (v NullableAccountGetRes) IsSet() bool {
	return v.isSet
}

func (v *NullableAccountGetRes) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAccountGetRes(val *AccountGetRes) *NullableAccountGetRes {
	return &NullableAccountGetRes{value: val, isSet: true}
}

func (v NullableAccountGetRes) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAccountGetRes) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

