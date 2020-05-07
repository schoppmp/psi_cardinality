package server

/*
#include "psi_cardinality/c/psi_cardinality_server_c.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
)

//PSICardinalityServer context
type PSICardinalityServer struct {
	context C.psi_cardinality_server_ctx
}

//CreateWithNewKey returns a new PSI server
func CreateWithNewKey() (*PSICardinalityServer, error) {
	psiServer := new(PSICardinalityServer)

	var err *C.char
	rcode := C.psi_cardinality_server_create_with_new_key(&psiServer.context, &err)
	if rcode != 0 {
		return nil, fmt.Errorf("failed to create server context: %v(%v)", psiServer.loadCString(&err), rcode)
	}
	if psiServer.context == nil {
		return nil, errors.New("failed to create server context: null")
	}

	runtime.SetFinalizer(psiServer, func(s *PSICardinalityServer) { s.destroy() })
	return psiServer, nil
}

//CreateFromKey returns a new PSI server
func CreateFromKey(key string) (*PSICardinalityServer, error) {
	psiServer := new(PSICardinalityServer)

	var err *C.char
	rcode := C.psi_cardinality_server_create_from_key(C.struct_server_buffer_t{
		buff:     C.CString(key),
		buff_len: C.ulong(len(key)),
	},
		&psiServer.context, &err)

	if rcode != 0 {
		return nil, fmt.Errorf("failed to create server context: %v(%v)", psiServer.loadCString(&err), rcode)
	}
	if psiServer.context == nil {
		return nil, errors.New("failed to create server context: null")
	}

	runtime.SetFinalizer(psiServer, func(s *PSICardinalityServer) { s.destroy() })
	return psiServer, nil
}

//CreateSetupMessage returns a message from server's dataset to be sent to the client
func (s *PSICardinalityServer) CreateSetupMessage(fpr float64, inputCount int64, rawInput []string) (string, error) {
	if s.context == nil {
		return "", errors.New("invalid context")
	}

	input := []C.struct_server_buffer_t{}
	for idx := range rawInput {
		input = append(input, C.struct_server_buffer_t{
			buff:     C.CString(rawInput[idx]),
			buff_len: C.ulong(len(rawInput[idx])),
		})
	}

	var out *C.char
	var err *C.char
	var outlen C.ulong

	rcode := C.psi_cardinality_server_create_setup_message(s.context, C.double(fpr), C.long(inputCount), &input[0], C.ulong(len(input)), &out, &outlen, &err)

	if rcode != 0 {
		return "", fmt.Errorf("setup_message failed: %v(%v)", s.loadCString(&err), rcode)
	}
	result := s.loadCString(&out)

	return result, nil
}

//ProcessRequest processes the client request
func (s *PSICardinalityServer) ProcessRequest(request string) (string, error) {
	if s.context == nil {
		return "", errors.New("invalid context")
	}

	var out *C.char
	var err *C.char
	var outlen C.ulong

	rcode := C.psi_cardinality_server_process_request(s.context, C.struct_server_buffer_t{
		buff:     C.CString(request),
		buff_len: C.ulong(len(request)),
	}, &out, &outlen, &err)

	if rcode != 0 {
		return "", fmt.Errorf("process request failed: %v(%v)", s.loadCString(&err), rcode)
	}

	result := s.loadCString(&out)
	return result, nil

}

//GetPrivateKeyBytes returns this instance's private key
func (s *PSICardinalityServer) GetPrivateKeyBytes() (string, error) {
	if s.context == nil {
		return "", errors.New("invalid context")
	}

	var out *C.char
	var outlen C.ulong
	var err *C.char

	rcode := C.psi_cardinality_server_get_private_key_bytes(s.context, &out, &outlen, &err)

	if rcode != 0 {
		return "", fmt.Errorf("get private keys failed: %v(%v)", s.loadCString(&err), rcode)
	}

	result := s.loadCString(&out)

	return result, nil
}

func (s *PSICardinalityServer) destroy() error {
	if s.context == nil {
		return errors.New("invalid context")
	}
	C.psi_cardinality_server_delete(&s.context)
	s.context = nil
	return nil
}

func (c *PSICardinalityServer) loadCString(buff **C.char) string {
	str := C.GoString(*buff)
	C.psi_cardinality_server_delete_buffer(c.context, buff)
	return str
}
