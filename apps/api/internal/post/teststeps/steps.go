package teststeps

import (
	"context"
	"slices"
	"strconv"
	"time"

	"github.com/cucumber/godog"
	"github.com/vandad1901/p3s/apps/api/internal/post"
	"github.com/vandad1901/p3s/apps/api/internal/scenario"
	"github.com/vandad1901/p3s/packages/go/godogutil"
	"github.com/vandad1901/p3s/packages/go/idv"
)

func CreateStep(s *scenario.Scenario) func(context.Context, string, *godog.Table) {
	return func(ctx context.Context, expectError string, table *godog.Table) {
		currentMap := s.SyncTableToMap(table)

		for _, row := range currentMap {
			rowPost := resolvePost(s.SharedData, row)
			addedKeys, mutateRequest := resolvePostBlockMutateList(s.SharedData, row)

			idv, addedIDs, err := s.A.PostService.CreatePost(ctx, rowPost, mutateRequest.Added)
			if expectError != "" {
				s.Require.Error(err)
				s.ReturnedError = err

				continue
			} else {
				s.Require.NoError(err)
			}

			row["ID"] = strconv.FormatInt(idv.ID, 10)
			row["UpdatedAt"] = idv.UpdatedAt.Format(godogutil.TimeFMT)

			for i, item := range addedKeys {
				s.DataMap[item]["ID"] = strconv.FormatInt(addedIDs[i], 10)
				s.DataMap[item]["PostID"] = strconv.FormatInt(idv.ID, 10)
			}
		}
	}
}

func GetStep(s *scenario.Scenario) func(context.Context, string, *godog.Table) {
	return func(ctx context.Context, expectError string, table *godog.Table) {
		currentMap := s.SyncTableToMap(table)

		for _, row := range currentMap {
			rowPost := resolvePost(s.SharedData, row)
			postBlockKeys := godogutil.ResolveCommaSeparatedValues(s.SharedData, row, AllPostBlocksKey)
			rowPostBlocks := make([]*post.PostBlock, len(postBlockKeys))

			for i, key := range postBlockKeys {
				rowPostBlocks[i] = resolvePostBlock(s.SharedData, s.DataMap[key])
			}

			actualPost, actualPostBlocks, err := s.A.PostService.GetPost(ctx, rowPost.Slug)
			if expectError != "" {
				s.Require.Error(err)
				s.ReturnedError = err

				continue
			} else {
				s.Require.NoError(err)
			}

			assertPost(s, rowPost, actualPost)
			assertPostBlocks(s, rowPostBlocks, actualPostBlocks)
		}
	}
}

func GetPostsStep(s *scenario.Scenario) func(context.Context, string, *godog.Table) {
	return func(ctx context.Context, expectError string, table *godog.Table) {
		currentMap := s.SyncTableToMap(table)

		for _, row := range currentMap {
			actualPosts, _, err := s.A.PostService.GetPosts(ctx, post.PostsQuery{
				AuthorID: godogutil.ResolveInt64(s.SharedData, row, "AuthorID"),

				CreatedBefore: time.Now(),
			})
			if expectError != "" {
				s.Require.Error(err)
				s.ReturnedError = err

				continue
			} else {
				s.Require.NoError(err)
			}

			ExpectedPostsKeys := godogutil.ResolveCommaSeparatedValues(s.SharedData, row, "ExpectedPosts")
			MissingPostsKeys := godogutil.ResolveCommaSeparatedValues(s.SharedData, row, "MissingPosts")

			ExpectedPosts := make([]*post.Post, len(ExpectedPostsKeys))
			for i, key := range ExpectedPostsKeys {
				ExpectedPosts[i] = resolvePost(s.SharedData, s.DataMap[key])
			}

			MissingPosts := make([]*post.Post, len(MissingPostsKeys))
			for i, key := range MissingPostsKeys {
				MissingPosts[i] = resolvePost(s.SharedData, s.DataMap[key])
			}

			for _, expected := range ExpectedPosts {
				actualIndex := slices.IndexFunc(actualPosts, func(candidate *post.Post) bool {
					return expected.ID == candidate.ID
				})

				s.Require.NotEqual(-1, actualIndex)
			}

			for _, missing := range MissingPosts {
				actualIndex := slices.IndexFunc(actualPosts, func(candidate *post.Post) bool {
					return missing.ID == candidate.ID
				})

				s.Require.Equal(-1, actualIndex)
			}
		}
	}
}

func assertPost(s *scenario.Scenario, expected, actual *post.Post) {
	s.Require.Equal(expected.Title, actual.Title)
}

func assertPostBlocks(s *scenario.Scenario, expectedList, actualList []*post.PostBlock) {
	s.Require.Len(actualList, len(expectedList))

	for _, expected := range expectedList {
		actualIndex := slices.IndexFunc(actualList, func(candidate *post.PostBlock) bool {
			return expected.ID == candidate.ID
		})

		s.Require.NotEqual(-1, actualIndex)

		actual := actualList[actualIndex]

		s.Require.Equal(expected.Position, actual.Position)
		s.Require.Equal(expected.BlockType, actual.BlockType)
		s.Require.Equal(expected.MediaContent, actual.MediaContent)
		s.Require.Equal(expected.TextContent, actual.TextContent)
		s.Require.Equal(expected.Metadata, actual.Metadata)
	}
}

func MustBeDeletedStep(s *scenario.Scenario) func(context.Context, *godog.Table) {
	return func(ctx context.Context, table *godog.Table) {
		currentMap := s.SyncTableToMap(table)

		for _, row := range currentMap {
			rowPost := resolvePost(s.SharedData, row)

			_, _, err := s.A.PostService.GetPost(ctx, rowPost.Slug)
			s.Require.ErrorIs(err, post.ErrPostNotFount)
		}
	}
}

func UpdateStep(s *scenario.Scenario) func(context.Context, string, *godog.Table) {
	return func(ctx context.Context, expectError string, table *godog.Table) {
		currentMap := s.SyncTableToMap(table)

		for _, row := range currentMap {
			rowPost := resolvePost(s.SharedData, row)
			addedKeys, mutateRequest := resolvePostBlockMutateList(s.SharedData, row)

			idv, addedIDs, err := s.A.PostService.UpdatePost(ctx, rowPost, mutateRequest)
			if expectError != "" {
				s.Require.Error(err)
				s.ReturnedError = err

				continue
			} else {
				s.Require.NoError(err)
			}

			row["UpdatedAt"] = idv.UpdatedAt.Format(godogutil.TimeFMT)

			for i, item := range addedKeys {
				s.DataMap[item]["ID"] = strconv.FormatInt(addedIDs[i], 10)
				s.DataMap[item]["PostID"] = strconv.FormatInt(idv.ID, 10)
			}
		}
	}
}

func DeleteStep(s *scenario.Scenario) func(context.Context, string, *godog.Table) {
	return func(ctx context.Context, expectError string, table *godog.Table) {
		currentMap := s.SyncTableToMap(table)

		for _, row := range currentMap {
			rowPost := resolvePost(s.SharedData, row)

			err := s.A.PostService.DeletePost(ctx, &idv.IDV{ID: rowPost.ID, UpdatedAt: rowPost.UpdatedAt})
			if expectError != "" {
				s.Require.Error(err)
				s.ReturnedError = err

				continue
			} else {
				s.Require.NoError(err)
			}
		}
	}
}
