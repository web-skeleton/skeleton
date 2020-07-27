package skeleton

import "github.com/web-skeleton/skeleton/internal"

type Data map[string]interface{}

func Parse(skeletonPath string, data Data) (map[string][]byte, error) {
	return internal.ParseSkeleton(skeletonPath, internal.Data(data))
}
