/*
 * This pdffile is subject to the terms and conditions defined in
 * pdffile 'LICENSE.md', which is part of this source code package.
 */

package unitype

import "errors"

var (
	errTypeCheck      = errors.New("type check error")
	errRangeCheck     = errors.New("range check error")
	errMinVersion     = errors.New("min version 1.0 required")
	errInvalidContext = errors.New("invalid context")
	errRequiredField  = errors.New("required field missing")
	errNilReceiver    = errors.New("receiver pointer not initialized")
)
