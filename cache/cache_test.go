package cache

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type object struct {
	Prop      string `json:"prop"`
	OtherProp int64  `json:"other_prop"`
}

func Test_Success(t *testing.T) {
	s, err := New(2, 16)
	require.NoError(t, err)

	o := object{
		Prop:      "foo",
		OtherProp: 123,
	}
	encoded, err := json.Marshal(o)
	require.NoError(t, err)
	err = s.Set("foo_key", encoded)
	require.NoError(t, err)

	result, err := s.Get("foo_key")
	require.NoError(t, err)

	var decoded object
	err = json.Unmarshal(result, &decoded)

	assert.Equal(t, o, decoded)
}

func Test_Eviction_MultiShard(t *testing.T) {
	s, err := New(2, 16)
	require.NoError(t, err)

	o := object{
		Prop:      "foo",
		OtherProp: 123,
	}
	encoded, err := json.Marshal(o)
	require.NoError(t, err)
	err = s.Set("foo_key", encoded)
	require.NoError(t, err)

	time.Sleep(3 * time.Second)
	err = s.Set("bar_key", encoded)
	require.NoError(t, err)

	/*
		Item life has ended but no other item has been added to its shard, so it will still be returned
		when requested. Item `bar_key` is likely added in another shard of the cache (due to default bigcache hasher)
		as a result the `foo_key` is still available.
	*/

	result, err := s.Get("foo_key")
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func Test_Eviction_SingleShard(t *testing.T) {
	s, err := New(2, 1)
	require.NoError(t, err)

	o := object{
		Prop:      "foo",
		OtherProp: 123,
	}
	encoded, err := json.Marshal(o)
	require.NoError(t, err)
	err = s.Set("foo_key", encoded)
	require.NoError(t, err)

	time.Sleep(3 * time.Second)
	err = s.Set("bar_key", encoded)
	require.NoError(t, err)

	/*
		Item life has ended and another key has been added to the same shard (single shard cache in this case)
		as a result the "dead" item gets evicted and not return in the subsequent call
	*/

	result, err := s.Get("foo_key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Nil(t, result)
}
