package container

import "github.com/docker/docker/pkg/archive"

func prepareArchive(src, dst string, excludes []string) (*archive.TempArchive, error) {
	reader, err := archive.TarWithOptions(src, &archive.TarOptions{
		Compression:     archive.Uncompressed,
		ExcludePatterns: excludes,
	})
	if err != nil {
		return nil, err
	}
	err = archive.Unpack(reader, dst, &archive.TarOptions{})
	if err != nil {
		return nil, err
	}
	reader, err = archive.Tar(dst, archive.Gzip)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return archive.NewTempArchive(reader, "")
}
