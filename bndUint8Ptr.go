// Copyright 2014 Rana Ian. All rights reserved.
// Use of this source code is governed by The MIT License
// found in the accompanying LICENSE file.

package ora

/*
#include <oci.h>
#include "version.h"
*/
import "C"
import (
	"unsafe"
)

type bndUint8Ptr struct {
	stmt      *Stmt
	ocibnd    *C.OCIBind
	ociNumber C.OCINumber
	isNull    C.sb2
	value     *uint8
}

func (bnd *bndUint8Ptr) bind(value *uint8, position int, stmt *Stmt) error {
	bnd.stmt = stmt
	bnd.value = value
	if value == nil {
		bnd.isNull = C.sb2(-1)
	} else {
		r := C.OCINumberFromInt(
			bnd.stmt.ses.srv.env.ocierr, //OCIError            *err,
			unsafe.Pointer(value),       //const void          *inum,
			1, //uword               inum_length,
			C.OCI_NUMBER_UNSIGNED, //uword               inum_s_flag,
			&bnd.ociNumber)        //OCINumber           *number );
		if r == C.OCI_ERROR {
			return bnd.stmt.ses.srv.env.ociError()
		}
	}
	r := C.OCIBINDBYPOS(
		bnd.stmt.ocistmt,                  //OCIStmt      *stmtp,
		(**C.OCIBind)(&bnd.ocibnd),        //OCIBind      **bindpp,
		bnd.stmt.ses.srv.env.ocierr,       //OCIError     *errhp,
		C.ub4(position),                   //ub4          position,
		unsafe.Pointer(&bnd.ociNumber),    //void         *valuep,
		C.LENGTH_TYPE(C.sizeof_OCINumber), //sb8          value_sz,
		C.SQLT_VNU,                        //ub2          dty,
		unsafe.Pointer(&bnd.isNull),       //void         *indp,
		nil,           //ub2          *alenp,
		nil,           //ub2          *rcodep,
		0,             //ub4          maxarr_len,
		nil,           //ub4          *curelep,
		C.OCI_DEFAULT) //ub4          mode );
	if r == C.OCI_ERROR {
		return bnd.stmt.ses.srv.env.ociError()
	}
	return nil
}

func (bnd *bndUint8Ptr) setPtr() error {
	if bnd.isNull > C.sb2(-1) {
		r := C.OCINumberToInt(
			bnd.stmt.ses.srv.env.ocierr, //OCIError              *err,
			&bnd.ociNumber,              //const OCINumber       *number,
			C.uword(1),                  //uword                 rsl_length,
			C.OCI_NUMBER_UNSIGNED,       //uword                 rsl_flag,
			unsafe.Pointer(bnd.value))   //void                  *rsl );
		if r == C.OCI_ERROR {
			return bnd.stmt.ses.srv.env.ociError()
		}
	}
	return nil
}

func (bnd *bndUint8Ptr) close() (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = errR(value)
		}
	}()

	stmt := bnd.stmt
	bnd.stmt = nil
	bnd.ocibnd = nil
	bnd.value = nil
	stmt.putBnd(bndIdxUint8Ptr, bnd)
	return nil
}
