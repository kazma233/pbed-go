package bed

type Bed interface {
	UploadByPath(filePath string) (string, error)

	UploadByBytes(bs []byte, fname string) (string, error)
}
