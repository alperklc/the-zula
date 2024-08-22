import React from "react";
import { useParams, useSearchParams } from "react-router-dom";
import Layout, { styles as layoutStyles } from "../../../components/layout";
import List from "../../../components/list";
import { useAuth } from "../../../contexts/authContext";
import { Api, NoteChange, PaginationMeta } from "../../../types/Api.ts";
import { useUI } from "../../../contexts/uiContext";
import Breadcrumbs from "../../../components/breadcrumbs";
import { Query, QueryContext, QueryContextProvider } from "../../../contexts/queryContext";
import { filterEmptyValues } from "../../../utils/filter"

function NotesChangesList() {
  const { isMobile } = useUI()
  const { shortId } = useParams()
  const { query, changePageSize, paginate } = React.useContext(QueryContext)
  const [loading, setLoading] = React.useState(true);
  const [data, setData] = React.useState<{ meta?: PaginationMeta, items?: NoteChange[] }>();
  const [error, setError] = React.useState<string | null>(null);
  
  const { user } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}` } } })

  const fetchAndUpdate = async (query: Query) => {
    try {
      setLoading(true);
      setError(null);

      const filteredQuery = filterEmptyValues(query)
      const { data, error, status } = await api.api.getNotesChanges(shortId!, filteredQuery)

      if (status === 200) {
        setData(data);
      } else {
        console.error(error);
        setError(error);
      }

    } catch (e: unknown) {
      console.error(e);
      setError(e as string);
    }
    setLoading(false);
  };

  React.useEffect(() => {
    const handler = setTimeout(
      () => fetchAndUpdate(query),
      100
    );

    return () => clearTimeout(handler);
  }, [query.page, query.pageSize]);

  return (
    <Layout
      fixedSubHeader={!isMobile}
      subHeaderContent={
        <>
          <Breadcrumbs />
          {!isMobile && (
            <div className={layoutStyles.subheader}>
              <div className={layoutStyles.subheaderTitleContainer}>
              </div>
              <></>
            </div>
          )}
        </>
      }
    >
      <List<NoteChange>
        resourceType="note-change"
        loading={loading}
        error={error}
        meta={data?.meta}
        items={data?.items}
        changePageSize={changePageSize}
        paginate={paginate}
      />
    </Layout>
  );
}

export function NotesChangesListPage() {
  const [params] = useSearchParams();
  
  const query = {
    tags: params.get('tags') || '',
    q: params.get('q') || '',
    page: parseInt(params.get('page') as string) || 1,
    pageSize: parseInt(params.get('pageSize') as string) || 10,
    sortBy: params.get('sortBy') || 'updatedAt',
    sortDirection: params.get('sortDirection') || 'desc',
  }

  return (
    <QueryContextProvider query={query}>
      <NotesChangesList />
    </QueryContextProvider>
  )
}

