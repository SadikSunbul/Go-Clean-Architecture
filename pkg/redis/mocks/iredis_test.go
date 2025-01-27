package mocks

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewRedis(t *testing.T) {
	// Arrange
	mockRedis := NewIRedis(t)

	// Mock davranışını ayarla
	mockRedis.On("IsConnected").Return(true)

	// Assert
	assert.NotEmpty(t, mockRedis)
	assert.NotNil(t, mockRedis)
	assert.True(t, mockRedis.IsConnected())
}

func Test_Redis_Get(t *testing.T) {
	// Arrange
	mockRedis := NewIRedis(t)

	// Mock davranışını ayarla
	mockRedis.On("Get", "key", "value").Return(nil)

	// Action
	err := mockRedis.Get("key", "value")

	// Assert
	assert.Nil(t, err)
}

func Test_Redis_Set(t *testing.T) {
	// Arrange
	mockRedis := NewIRedis(t)
	mockRedis.On("Set", "key", "value").Return(nil)

	// Action
	err := mockRedis.Set("key", "value")

	// Assert
	assert.Nil(t, err)
	mockRedis.AssertExpectations(t)
}

func Test_Redis_SetWithExpiration(t *testing.T) {
	// Arrange
	mockRedis := NewIRedis(t)
	duration := time.Hour * 1
	mockRedis.On("SetWithExpiration", "key", "value", duration).Return(nil)

	// Action
	err := mockRedis.SetWithExpiration("key", "value", duration)

	// Assert
	assert.Nil(t, err)
	mockRedis.AssertExpectations(t)
}

func Test_Redis_Remove(t *testing.T) {
	// Arrange
	mockRedis := NewIRedis(t)
	mockRedis.On("Remove", "key").Return(nil)

	// Action
	err := mockRedis.Remove("key")

	// Assert
	assert.Nil(t, err)
	mockRedis.AssertExpectations(t)
}

func Test_Redis_RemovePattern(t *testing.T) {
	// Arrange
	mockRedis := NewIRedis(t)
	mockRedis.On("RemovePattern", "key*").Return(nil)

	// Action
	err := mockRedis.RemovePattern("key*")

	// Assert
	assert.Nil(t, err)
	mockRedis.AssertExpectations(t)
}

func Test_Redis_Get_NotFound(t *testing.T) {
	// Arrange
	mockRedis := NewIRedis(t)
	var value string
	mockRedis.On("Get", "nonexistent", &value).Return(errors.New("key not found"))

	// Action
	err := mockRedis.Get("nonexistent", &value)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key not found")
	mockRedis.AssertExpectations(t)
}

func Test_Redis_Set_Error(t *testing.T) {
	// Arrange
	mockRedis := NewIRedis(t)
	mockRedis.On("Set", "key", "value").Return(errors.New("connection error"))

	// Action
	err := mockRedis.Set("key", "value")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection error")
	mockRedis.AssertExpectations(t)
}

func Test_Redis_IsConnected_False(t *testing.T) {
	// Arrange
	mockRedis := NewIRedis(t)
	mockRedis.On("IsConnected").Return(false)

	// Assert
	assert.False(t, mockRedis.IsConnected())
	mockRedis.AssertExpectations(t)
}

func Test_Redis_GetWithStruct(t *testing.T) {
	// Arrange
	mockRedis := NewIRedis(t)
	type TestStruct struct {
		Name string
		Age  int
	}
	value := TestStruct{Name: "test", Age: 25}
	mockRedis.On("Get", "struct-key", &value).Return(nil)

	// Action
	err := mockRedis.Get("struct-key", &value)

	// Assert
	assert.Nil(t, err)
	mockRedis.AssertExpectations(t)
}

func Test_Redis_SetWithStruct(t *testing.T) {
	// Arrange
	mockRedis := NewIRedis(t)
	value := struct {
		Name string
		Age  int
	}{
		Name: "test",
		Age:  25,
	}
	mockRedis.On("Set", "struct-key", value).Return(nil)

	// Action
	err := mockRedis.Set("struct-key", value)

	// Assert
	assert.Nil(t, err)
	mockRedis.AssertExpectations(t)
}
