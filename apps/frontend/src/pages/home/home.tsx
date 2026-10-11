import { postService } from "@/api/api.service";
import { Button } from "@/components/ui/button";
import { type Post } from "@gen/api/postpb/v1/post";
import { ChevronLeft, ChevronRight, Loader2 } from "lucide-react";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

export function Home() {
  const [loadTime, _] = useState(new Date().toISOString());

  const [page, setPage] = useState(0);
  const perPage = 10;
  const [posts, setPosts] = useState<Post[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    setIsLoading(true);
    postService
      .GetPosts({ authorId: "0", limit: perPage, offset: page * perPage, createdBefore: loadTime })
      .then((res) => {
        if (res.ok) {
          setPosts(res.posts);
        } else {
          console.error(res.message);
        }
      })
      .finally(() => setIsLoading(false));
  }, [page, perPage]);

  return (
    <div className="mx-auto flex w-full max-w-3xl flex-1 flex-col items-center justify-center p-4 py-2 md:p-0">
      {isLoading ? (
        <div className="text-muted-foreground flex flex-col items-center gap-1">
          <p>Loading</p>
          <Loader2 className="animate-spin" />
        </div>
      ) : (
        <div className="flex h-full w-full flex-1 flex-col items-center gap-2">
          <h1 className="w-full border-b pt-1 pb-3 text-3xl">All Posts</h1>
          <ul className="flex w-full max-w-2xl flex-1 flex-col">
            {posts.map((post) => (
              <li
                key={post.id}
                className="flex justify-between border-b border-gray-300 py-4"
              >
                <Link to={`/post/${post.slug}`}>{post.title}</Link>
                <p className="text-muted-foreground">{new Date(post.createdAt!).toDateString()}</p>
              </li>
            ))}
          </ul>

          <div className="flex justify-center gap-6 py-2">
            <Button
              variant={"outline"}
              onClick={() => setPage(page - 1)}
              disabled={page === 0}
              className="flex"
            >
              <ChevronLeft />
              Prev
            </Button>
            ...
            <Button
              variant={"outline"}
              onClick={() => setPage(page + 1)}
              disabled={posts.length < perPage}
              className="flex"
            >
              Next
              <ChevronRight />
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
