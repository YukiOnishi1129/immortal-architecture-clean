import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { listMyNoteQuery } from "@/external/handler/note/note.query.server";
import { MyNoteList } from "@/features/note/components/client/MyNoteList";
import { noteKeys } from "@/features/note/queries/keys";
import type { NoteStatus } from "@/features/note/types";
import { getQueryClient } from "@/shared/lib/query-client";

type MyNoteListPageTemplateProps = {
  status?: NoteStatus;
  q?: string;
  page?: number;
};

export async function MyNoteListPageTemplate(
  props: MyNoteListPageTemplateProps = {},
) {
  const queryClient = getQueryClient();

  const filters = {
    status: props.status,
    q: props.q,
    page: props.page,
  };

  // ownerIdは認証情報から自動的に設定されるため、ここでは設定不要
  await queryClient.prefetchQuery({
    queryKey: noteKeys.myList(filters),
    queryFn: () => listMyNoteQuery(filters),
  });

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <MyNoteList initialFilters={filters} />
    </HydrationBoundary>
  );
}
