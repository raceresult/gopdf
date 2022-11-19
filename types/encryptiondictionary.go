package types

// PDF Reference 1.4, Table 3.13 Entries common to all encryption dictionaries

type EncryptionDictionary struct {
	// (Required) The name of the security handler for this document; see below. Default value:
	// Standard, for the built-in security handler. (Names for other security handlers can be
	// registered using the procedure described in Appendix E.)
	Filter Name

	// (Optional but strongly recommended) A code specifying the algorithm to be used in en-
	// crypting and decrypting the document:
	// 0 An algorithm that is undocumented and no longer supported, and whose use is
	// strongly discouraged.
	// 1 Algorithm 3.1 on page 73, with an encryption key length of 40 bits; see below.
	// 2 (PDF 1.4) Algorithm 3.1 on page 73, but allowing encryption key lengths greater
	// than 40 bits.
	// 3 (PDF 1.4) An unpublished algorithm allowing encryption key lengths ranging
	// from 40 to 128 bits. (This algorithm is unpublished as an export requirement of
	// the U.S. Department of Commerce.)
	// The default value if this entry is omitted is 0, but a value of 1 or greater is strongly rec-
	// ommended. (See implementation note 15 in Appendix H.)
	V Number

	// (Optional; PDF 1.4; only if V is 2 or 3) The length of the encryption key, in bits. The value
	// must be a multiple of 8, in the range 40 to 128. Default value: 40.
	Length Int
}

func (q EncryptionDictionary) ToRawBytes() []byte {
	d := Dictionary{
		"Filter": q.Filter,
		"V":      q.V,
	}
	if q.Length != 0 {
		d["Length"] = q.Length
	}
	return d.ToRawBytes()
}

func (q EncryptionDictionary) Copy(copyRef func(reference Reference) Reference) Object {
	return EncryptionDictionary{
		Filter: q.Filter.Copy(copyRef).(Name),
		V:      q.V.Copy(copyRef).(Number),
		Length: q.Length.Copy(copyRef).(Int),
	}
}

func (q EncryptionDictionary) Equal(obj Object) bool {
	a, ok := obj.(EncryptionDictionary)
	if !ok {
		return false
	}
	if !q.Filter.Equal(a.Filter) {
		return false
	}
	if !q.V.Equal(a.V) {
		return false
	}
	if !q.Length.Equal(a.Length) {
		return false
	}
	return true
}

// PDF Reference 1.4, Table 3.14 Additional encryption dictionary entries for the standard security handler

type StandardSecurityHandler struct {
	EncryptionDictionary

	// (Required) A number specifying which revision of the standard security handler should
	// be used to interpret this dictionary. The revision number should be 2 if the document is
	// encrypted with a V value less than 2 (see Table 3.13) and does not have any of the access
	// permissions set (via the P entry, below) that are designated “Revision 3” in Table 3.15;
	// otherwise (that is, if the document is encrypted with a V value greater than 2 or has any
	// “Revision 3” access permissions set), this value should be 3.
	R Number

	// (Required) A 32-byte string, based on both the owner and user passwords, that is used in
	// computing the encryption key and in determining whether a valid owner password was
	// entered. For more information, see “Encryption Key Algorithm” on page 78 and “Pass-
	// word Algorithms” on page 79.
	O String

	// (Required) A 32-byte string, based on the user password, that is used in determining
	// whether to prompt the user for a password and, if so, whether a valid user or owner pass-
	// word was entered. For more information, see “Password Algorithms” on page 79.
	U String

	// (Required) A set of flags specifying which operations are permitted when the document is
	// opened with user access (see Table 3.15).
	P Int
}

func (q StandardSecurityHandler) ToRawBytes() []byte {
	d := Dictionary{
		"Filter": q.Filter,
		"V":      q.V,
		"R":      q.R,
		"O":      q.O,
		"U":      q.U,
		"P":      q.P,
	}
	if q.Length != 0 {
		d["Length"] = q.Length
	}
	return d.ToRawBytes()
}

func (q StandardSecurityHandler) Copy(copyRef func(reference Reference) Reference) Object {
	return StandardSecurityHandler{
		EncryptionDictionary: Copy(q.EncryptionDictionary, copyRef).(EncryptionDictionary),
		R:                    q.R.Copy(copyRef).(Number),
		O:                    q.O.Copy(copyRef).(String),
		U:                    q.U.Copy(copyRef).(String),
		P:                    q.P.Copy(copyRef).(Int),
	}
}

func (q StandardSecurityHandler) Equal(obj Object) bool {
	a, ok := obj.(StandardSecurityHandler)
	if !ok {
		return false
	}
	if !Equal(q.EncryptionDictionary, a.EncryptionDictionary) {
		return false
	}
	if !Equal(q.R, a.R) {
		return false
	}
	if !Equal(q.O, a.O) {
		return false
	}
	if !Equal(q.U, a.U) {
		return false
	}
	if !Equal(q.P, a.P) {
		return false
	}
	return true
}
