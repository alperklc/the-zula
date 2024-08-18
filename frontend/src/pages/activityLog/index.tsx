import React from "react";
import { useSearchParams } from "react-router-dom";
import Layout, { styles as layoutStyles } from "../../components/layout";
import List from "../../components/list";
import { useAuth } from "../../contexts/authContext";
import { Api, PaginationMeta, UserActivity } from "../../types/Api.ts";
import { useUI } from "../../contexts/uiContext";
import Breadcrumbs from "../../components/breadcrumbs";
import { Query, QueryContext, QueryContextProvider } from "../../contexts/queryContext";
import { filterEmptyValues } from "../../utils/filter"

function ActivitiesList() {
  const { isMobile } = useUI()

  const { query, changePageSize, paginate } = React.useContext(QueryContext)
  const [loading, setLoading] = React.useState(true);
  const [data, setData] = React.useState<{ meta?: PaginationMeta, items?: UserActivity[] }>();
  const [error, setError] = React.useState<string | null>(null);
  
  const { user } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}` } } })

  const fetchAndUpdate = async (query: Query) => {
    try {
      setLoading(true);
      setError(null);

      const filteredQuery = filterEmptyValues(query)
      const { data, status } = await api.api.getUserActivity(`${user?.profile.sub}`, filteredQuery)

      if (status === 200) {
        setData(data);
      } else {
        console.error(data);
        setError(data);
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
  }, [query.page, query.pageSize, query.sortBy, query.sortDirection]);

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
      <List<UserActivity>
        resourceType="user-activity"
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

export function ActivitiesListPage() {
  const [params] = useSearchParams();
  
  const query = {
    tags: params.get('tags') || '',
    q: params.get('q') || '',
    page: parseInt(params.get('page') as string) || 1,
    pageSize: parseInt(params.get('pageSize') as string) || 10,
    sortBy: params.get('sortBy') || 'timestamp',
    sortDirection: params.get('sortDirection') || 'desc',
  }

  return (
    <QueryContextProvider query={query}>
      <ActivitiesList />
    </QueryContextProvider>
  )
}

