// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	dao "github.com/julkhong/walletapp/server/internal/dao"
	mock "github.com/stretchr/testify/mock"
)

// WalletDaoInterface is an autogenerated mock type for the WalletDaoInterface type
type WalletDaoInterface struct {
	mock.Mock
}

// CheckIdempotencyKey provides a mock function with given fields: key, method, path
func (_m *WalletDaoInterface) CheckIdempotencyKey(key string, method string, path string) (*dao.IdempotencyRecord, bool) {
	ret := _m.Called(key, method, path)

	if len(ret) == 0 {
		panic("no return value specified for CheckIdempotencyKey")
	}

	var r0 *dao.IdempotencyRecord
	var r1 bool
	if rf, ok := ret.Get(0).(func(string, string, string) (*dao.IdempotencyRecord, bool)); ok {
		return rf(key, method, path)
	}
	if rf, ok := ret.Get(0).(func(string, string, string) *dao.IdempotencyRecord); ok {
		r0 = rf(key, method, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.IdempotencyRecord)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, string) bool); ok {
		r1 = rf(key, method, path)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// CreateTransaction provides a mock function with given fields: tx
func (_m *WalletDaoInterface) CreateTransaction(tx *dao.Transaction) error {
	ret := _m.Called(tx)

	if len(ret) == 0 {
		panic("no return value specified for CreateTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*dao.Transaction) error); ok {
		r0 = rf(tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetBalance provides a mock function with given fields: walletID
func (_m *WalletDaoInterface) GetBalance(walletID string) (float64, error) {
	ret := _m.Called(walletID)

	if len(ret) == 0 {
		panic("no return value specified for GetBalance")
	}

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (float64, error)); ok {
		return rf(walletID)
	}
	if rf, ok := ret.Get(0).(func(string) float64); ok {
		r0 = rf(walletID)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(walletID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionHistory provides a mock function with given fields: walletID, txType, start, end, limit, offset
func (_m *WalletDaoInterface) GetTransactionHistory(walletID string, txType string, start string, end string, limit int, offset int) ([]dao.Transaction, error) {
	ret := _m.Called(walletID, txType, start, end, limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetTransactionHistory")
	}

	var r0 []dao.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, string, string, int, int) ([]dao.Transaction, error)); ok {
		return rf(walletID, txType, start, end, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(string, string, string, string, int, int) []dao.Transaction); ok {
		r0 = rf(walletID, txType, start, end, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dao.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, string, string, int, int) error); ok {
		r1 = rf(walletID, txType, start, end, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWalletByID provides a mock function with given fields: walletID
func (_m *WalletDaoInterface) GetWalletByID(walletID string) (*dao.Wallet, error) {
	ret := _m.Called(walletID)

	if len(ret) == 0 {
		panic("no return value specified for GetWalletByID")
	}

	var r0 *dao.Wallet
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*dao.Wallet, error)); ok {
		return rf(walletID)
	}
	if rf, ok := ret.Get(0).(func(string) *dao.Wallet); ok {
		r0 = rf(walletID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.Wallet)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(walletID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveIdempotencyKey provides a mock function with given fields: record
func (_m *WalletDaoInterface) SaveIdempotencyKey(record *dao.IdempotencyRecord) error {
	ret := _m.Called(record)

	if len(ret) == 0 {
		panic("no return value specified for SaveIdempotencyKey")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*dao.IdempotencyRecord) error); ok {
		r0 = rf(record)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateBalance provides a mock function with given fields: input
func (_m *WalletDaoInterface) UpdateBalance(input *dao.UpdateBalance) error {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for UpdateBalance")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*dao.UpdateBalance) error); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewWalletDaoInterface creates a new instance of WalletDaoInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWalletDaoInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *WalletDaoInterface {
	mock := &WalletDaoInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
