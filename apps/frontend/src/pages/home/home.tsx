import type { APIResponse } from "@/api/client";
import { postService } from "@/api/post.service";
import { userService } from "@/api/user.service";
import { type Post } from "@gen/api/postpb/v1/post";
import { Loader2 } from "lucide-react";
import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { PostsPagination } from "./pagination";

type EnrichedPost = {
  post: Post | null;
  authorName: string;
};

async function fetchPosts(
  perPage: number,
  page: number,
  loadTime: string,
): Promise<APIResponse<{ data: EnrichedPost[]; total: number }>> {
  debugger;
  const res = await postService.GetPosts({
    authorId: "0",
    limit: perPage,
    offset: (page - 1) * perPage,
    createdBefore: loadTime,
  });

  if (!res.ok) {
    return { ok: false, message: res.message, code: res.code };
  }

  const authorIds = res.posts.map((p) => p.createdBy);
  const usersRes = await userService.GetUsers({ userIds: authorIds });
  if (!usersRes.ok) {
    return { ok: false, message: usersRes.message, code: usersRes.code };
  }

  const enrichedPosts: EnrichedPost[] = res.posts.map((post) => {
    const author = usersRes.users.find((u) => u.id === post.createdBy);
    return {
      post,
      authorName: author ? author.displayName : "",
    };
  });

  return { ok: true, data: enrichedPosts, total: Number(res.total) };
}

export function Home() {
  const [loadTime] = useState(new Date().toISOString());

  const [searchParams] = useSearchParams();
  const perPage = 10;
  const page = Number(searchParams.get("page")) || 1;

  const [totalPosts, setTotalPosts] = useState(0);
  const totalPages = Math.ceil(totalPosts / perPage);

  const [posts, setPosts] = useState<EnrichedPost[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    setIsLoading(true);
    let cancelled = false;
    fetchPosts(perPage, page, loadTime)
      .then((res) => {
        if (cancelled) return;
        if (res.ok) {
          setPosts(res.data);
          setTotalPosts(res.total);
        } else {
          console.error(res.message);
        }
      })
      .finally(() => setIsLoading(false));
    return () => {
      cancelled = true;
    };
  }, [page, perPage]);

  return (
    <div className="mx-auto flex w-full max-w-3xl flex-1 flex-col items-center justify-center p-4 py-2 md:p-0">
      <div className="flex h-full w-full flex-1 flex-col items-center gap-1">
        <h1 className="w-full border-b pt-1 pb-3 text-3xl">All Posts</h1>
        <PostsPagination
          page={page}
          totalPages={totalPages}
        />
        {isLoading ? (
          <div className="text-muted-foreground flex flex-col items-center gap-1">
            <p>Loading</p>
            <Loader2 className="animate-spin" />
          </div>
        ) : (
          <ul className="flex w-full max-w-2xl flex-1 flex-col">
            {posts.map(({ post, authorName }) => (
              <li
                key={post!.id}
                className="flex justify-between gap-4 border-b border-gray-300 py-4"
              >
                <div className="flex flex-col">
                  <Link
                    to={`/post/${post!.slug}`}
                    className="text-lg break-all hyphens-auto"
                  >
                    {post!.title}
                  </Link>
                  <p className="text-muted-foreground text-md">{authorName}</p>
                </div>
                <p className="text-muted-foreground text-md shrink-0">
                  {new Date(post!.createdAt!).toDateString()}
                </p>
              </li>
            ))}
          </ul>
        )}
      </div>
      <PostsPagination
        page={page}
        totalPages={totalPages}
      />
    </div>
  );
}
