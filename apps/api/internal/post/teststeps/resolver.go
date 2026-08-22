package teststeps

import (
	"maps"
	"slices"
	"strings"

	"github.com/vandad1901/p3s/apps/api/internal/post"
	"github.com/vandad1901/p3s/packages/go/godogutil"
)

func resolvePost(s *godogutil.SharedData, row map[string]string) *post.Post {
	var postStatusEnums = [3]string{"unspecified", "draft", "published"}

	p := new(post.Post)

	p.ID = godogutil.ResolveInt64(s, row, "ID")

	p.Title = godogutil.ResolveString(s, row, "Title")
	p.Slug = godogutil.ResolveString(s, row, "Slug")
	p.Status = post.PostStatus(godogutil.ResolveEnumFactory(s, row, "Status", postStatusEnums[:]))

	p.CreatedAt = godogutil.ResolveTime(s, row, "CreatedAt")
	p.CreatedBy = godogutil.ResolveInt64(s, row, "CreatedBy")
	p.UpdatedAt = godogutil.ResolveTime(s, row, "UpdatedAt")
	p.UpdatedBy = godogutil.ResolveInt64(s, row, "UpdatedBy")

	return p
}

func resolvePostBlock(s *godogutil.SharedData, row map[string]string) *post.PostBlock {
	var blockTypeEnums = [3]string{"unspecified", "text", "media"}

	pb := new(post.PostBlock)

	pb.ID = godogutil.ResolveInt64(s, row, "ID")
	pb.PostID = godogutil.ResolveInt64(s, row, "PostID")
	pb.Position = godogutil.ResolveInt32(s, row, "Position")

	pb.BlockType = post.BlockType(godogutil.ResolveEnumFactory(s, row, "BlockType", blockTypeEnums[:]))
	pb.TextContent = godogutil.ResolveString(s, row, "TextContent")
	pb.MediaContent = godogutil.ResolveString(s, row, "MediaContent")
	pb.Metadata = godogutil.ResolveString(s, row, "Metadata")

	return pb
}

const (
	AllPostBlocksKey     = "AllPostBlocks"
	AddedPostBlocksKey   = "AddedPostBlocks"
	EditedPostBlocksKey  = "EditedPostBlocks"
	RemovedPostBlocksKey = "RemovedPostBlocks"
)

func resolvePostBlockMutateList(s *godogutil.SharedData, row map[string]string,
) ([]string, *post.PostBlockMutateRequest) {
	req := new(post.PostBlockMutateRequest)

	var (
		ExistingKeys = godogutil.ResolveCommaSeparatedValues(s, row, AllPostBlocksKey)
		AddedKeys    = godogutil.ResolveCommaSeparatedValues(s, row, AddedPostBlocksKey)
		EditedKeys   = godogutil.ResolveCommaSeparatedValues(s, row, EditedPostBlocksKey)
		RemovedKeys  = godogutil.ResolveCommaSeparatedValues(s, row, RemovedPostBlocksKey)
	)

	allKeysMap := make(map[string]struct{})
	for _, rowKey := range ExistingKeys {
		allKeysMap[rowKey] = struct{}{}
	}

	req.Added = make([]*post.PostBlock, len(AddedKeys))
	for i, rowKey := range AddedKeys {
		req.Added[i] = resolvePostBlock(s, s.DataMap[rowKey])
		allKeysMap[rowKey] = struct{}{}
	}

	req.Edited = make([]*post.PostBlock, len(EditedKeys))
	for i, rowKey := range EditedKeys {
		req.Edited[i] = resolvePostBlock(s, s.DataMap[rowKey])
	}

	req.Removed = make([]int64, len(RemovedKeys))
	for i, rowKey := range RemovedKeys {
		req.Removed[i] = godogutil.ResolveInt64(s, s.DataMap[rowKey], "ID")
		delete(allKeysMap, rowKey)
	}

	row[AddedPostBlocksKey] = ""
	row[EditedPostBlocksKey] = ""
	row[RemovedPostBlocksKey] = ""
	row[AllPostBlocksKey] = strings.Join(slices.Collect(maps.Keys(allKeysMap)), ",")

	return AddedKeys, req
}
