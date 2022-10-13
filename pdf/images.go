package pdf

import "git.rrdc.de/lib/gopdf/types"

// Image holds both the Image object and the reference to it
type Image struct {
	Reference types.Reference
	Image     *types.Image
}
