package types

import (
	"encoding/json"
	"github.com/guodongq/quickstart/pkg/util"
	"time"
)

var EmptyMeta = Meta{}

type Meta struct {
	CreatedBy *string    `bson:"created_by,omitempty"`
	CreatedAt *time.Time `bson:"created_at,omitempty"`

	UpdatedBy *string    `bson:"updated_by,omitempty"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`

	DeletedBy *string    `bson:"deleted_by,omitempty"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty"`

	Version *int32 `bson:"version,omitempty"`
}

func NewMeta() *Meta {
	this := Meta{}
	return &this
}

func NewMetaWithDefaults() *Meta {
	this := Meta{}
	return &this
}

// GetCreatedBy returns the CreatedBy field value if set, zero value otherwise.
func (o *Meta) GetCreatedBy() string {
	if o == nil || util.IsNil(o.CreatedBy) {
		var ret string
		return ret
	}
	return *o.CreatedBy
}

// GetCreatedByOk returns a tuple with the CreatedBy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Meta) GetCreatedByOk() (*string, bool) {
	if o == nil || util.IsNil(o.CreatedBy) {
		return nil, false
	}
	return o.CreatedBy, true
}

// HasCreatedBy returns a boolean if a field has been set.
func (o *Meta) HasCreatedBy() bool {
	if o != nil && !util.IsNil(o.CreatedBy) {
		return true
	}

	return false
}

// SetCreatedBy gets a reference to the given string and assigns it to the CreatedBy field.
func (o *Meta) SetCreatedBy(v string) {
	o.CreatedBy = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise.
func (o *Meta) GetCreatedAt() time.Time {
	if o == nil || util.IsNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Meta) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || util.IsNil(o.CreatedAt) {
		return nil, false
	}
	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *Meta) HasCreatedAt() bool {
	if o != nil && !util.IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *Meta) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetUpdatedBy returns the UpdatedBy field value if set, zero value otherwise.
func (o *Meta) GetUpdatedBy() string {
	if o == nil || util.IsNil(o.UpdatedBy) {
		var ret string
		return ret
	}
	return *o.UpdatedBy
}

// GetUpdatedByOk returns a tuple with the UpdatedBy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Meta) GetUpdatedByOk() (*string, bool) {
	if o == nil || util.IsNil(o.UpdatedBy) {
		return nil, false
	}
	return o.UpdatedBy, true
}

// HasUpdatedBy returns a boolean if a field has been set.
func (o *Meta) HasUpdatedBy() bool {
	if o != nil && !util.IsNil(o.UpdatedBy) {
		return true
	}

	return false
}

// SetUpdatedBy gets a reference to the given string and assigns it to the UpdatedBy field.
func (o *Meta) SetUpdatedBy(v string) {
	o.UpdatedBy = &v
}

// GetUpdatedAt returns the UpdatedAt field value if set, zero value otherwise.
func (o *Meta) GetUpdatedAt() time.Time {
	if o == nil || util.IsNil(o.UpdatedAt) {
		var ret time.Time
		return ret
	}
	return *o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Meta) GetUpdatedAtOk() (*time.Time, bool) {
	if o == nil || util.IsNil(o.UpdatedAt) {
		return nil, false
	}
	return o.UpdatedAt, true
}

// HasUpdatedAt returns a boolean if a field has been set.
func (o *Meta) HasUpdatedAt() bool {
	if o != nil && !util.IsNil(o.UpdatedAt) {
		return true
	}

	return false
}

// SetUpdatedAt gets a reference to the given time.Time and assigns it to the UpdatedAt field.
func (o *Meta) SetUpdatedAt(v time.Time) {
	o.UpdatedAt = &v
}

// GetDeletedBy returns the DeletedBy field value if set, zero value otherwise.
func (o *Meta) GetDeletedBy() string {
	if o == nil || util.IsNil(o.DeletedBy) {
		var ret string
		return ret
	}
	return *o.DeletedBy
}

// GetDeletedByOk returns a tuple with the DeletedBy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Meta) GetDeletedByOk() (*string, bool) {
	if o == nil || util.IsNil(o.DeletedBy) {
		return nil, false
	}
	return o.DeletedBy, true
}

// HasDeletedBy returns a boolean if a field has been set.
func (o *Meta) HasDeletedBy() bool {
	if o != nil && !util.IsNil(o.DeletedBy) {
		return true
	}

	return false
}

// SetDeletedBy gets a reference to the given string and assigns it to the DeletedBy field.
func (o *Meta) SetDeletedBy(v string) {
	o.DeletedBy = &v
}

// GetDeletedAt returns the DeletedAt field value if set, zero value otherwise.
func (o *Meta) GetDeletedAt() time.Time {
	if o == nil || util.IsNil(o.DeletedAt) {
		var ret time.Time
		return ret
	}
	return *o.DeletedAt
}

// GetDeletedAtOk returns a tuple with the DeletedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Meta) GetDeletedAtOk() (*time.Time, bool) {
	if o == nil || util.IsNil(o.DeletedAt) {
		return nil, false
	}
	return o.DeletedAt, true
}

// HasDeletedAt returns a boolean if a field has been set.
func (o *Meta) HasDeletedAt() bool {
	if o != nil && !util.IsNil(o.DeletedAt) {
		return true
	}

	return false
}

// SetDeletedAt gets a reference to the given time.Time and assigns it to the DeletedAt field.
func (o *Meta) SetDeletedAt(v time.Time) {
	o.DeletedAt = &v
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *Meta) GetVersion() int32 {
	if o == nil || util.IsNil(o.Version) {
		var ret int32
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Meta) GetVersionOk() (*int32, bool) {
	if o == nil || util.IsNil(o.Version) {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *Meta) HasVersion() bool {
	if o != nil && !util.IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given int32 and assigns it to the Version field.
func (o *Meta) SetVersion(v int32) {
	o.Version = &v
}

func (o *Meta) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o *Meta) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !util.IsNil(o.CreatedBy) {
		toSerialize["created_by"] = o.CreatedBy
	}
	if !util.IsNil(o.CreatedAt) {
		toSerialize["created_at"] = o.CreatedAt
	}
	if !util.IsNil(o.UpdatedBy) {
		toSerialize["updated_by"] = o.UpdatedBy
	}
	if !util.IsNil(o.UpdatedAt) {
		toSerialize["updated_at"] = o.UpdatedAt
	}
	if !util.IsNil(o.DeletedBy) {
		toSerialize["deleted_by"] = o.DeletedBy
	}
	if !util.IsNil(o.DeletedAt) {
		toSerialize["deleted_at"] = o.DeletedAt
	}
	if !util.IsNil(o.Version) {
		toSerialize["version"] = o.Version
	}
	return toSerialize, nil
}

type NullableMeta struct {
	value *Meta
	isSet bool
}

func (v *NullableMeta) Get() *Meta {
	return v.value
}

func (v *NullableMeta) Set(val *Meta) {
	v.value = val
	v.isSet = true
}

func (v *NullableMeta) IsSet() bool {
	return v.isSet
}

func (v *NullableMeta) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAudit(val *Meta) *NullableMeta {
	return &NullableMeta{value: val, isSet: true}
}

func (v *NullableMeta) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMeta) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
