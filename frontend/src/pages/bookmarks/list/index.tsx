import React from "react";
import { Link, useSearchParams } from "react-router-dom";
import { FormattedMessage } from "react-intl";
import Layout, { styles as layoutStyles } from "../../../components/layout";
import List from "../../../components/list";
import SearchInput from "../../../components/search";
import { useAuth } from "../../../contexts/authContext";
import { Api, Bookmark, PaginationMeta } from "../../../types/Api.ts";
import { useUI } from "../../../contexts/uiContext";
import Breadcrumbs from "../../../components/breadcrumbs";
import Button from "../../../components/form/button";
import { Query, QueryContext, QueryContextProvider } from "../../../contexts/queryContext";
import { filterEmptyValues } from "../../../utils/filter"
import MobileHeader from "../../../components/mobileHeader";
import useModal from "../../../components/modal";
import FilterSelector, { FilterSelectorModalProps } from "../../../components/filterSelectorModal";
import Icons from "../../../components/icons";

function BookmarksList() {
  const { isMobile } = useUI()
  const [FilterSelectionModal, openFilterSelection] =
    useModal<FilterSelectorModalProps>(FilterSelector)

  const { query, applyFilterAndSort, tagsAsArray, changePageSize, updateQuery, paginate } = React.useContext(QueryContext)
  const [loading, setLoading] = React.useState(true);
  const [data, setData] = React.useState<{ meta?: PaginationMeta, items?: Bookmark[] }>();
  const [error, setError] = React.useState<string | null>(null);
  
  const { user } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}` } } })

  const fetchAndUpdate = async (query: Query) => {
    try {
      setLoading(true);
      setError(null);

      const filteredQuery = filterEmptyValues(query)
      const { data, error, status } = await api.api.getBookmarks(filteredQuery)

      if (status === 200) {
        setData(data);
      } else {
        console.error(data);
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
  }, [query.page, query.pageSize, query.q, query.sortBy, query.sortDirection]);

  return (
    <Layout
      customHeader={
        isMobile
          ? (props: object) => <MobileHeader onFilterSelectionOpen={openFilterSelection} linkTo='/bookmarks/create' {...props}/>
          : undefined
      }
      fixedSubHeader={!isMobile}
      subHeaderContent={
        <>
          <Breadcrumbs />
          {!isMobile && (
            <div className={layoutStyles.subheader}>
              <div className={layoutStyles.subheaderTitleContainer}>
                <div className={layoutStyles.searchWithFilter}>
                  <SearchInput autoFocus query={query.q} onQueryUpdate={updateQuery} />
                  <Button onClick={openFilterSelection} className={layoutStyles.filterButton}>
                    <Icons.Filter width={18} height={18} />
                    {tagsAsArray.length > 0 && (
                      <span className={layoutStyles.badgeOnButton}>{tagsAsArray.length}</span>
                    )}
                  </Button>
                </div>
              </div>

              <Link to='/bookmarks/create'>
                <Button>
                  <FormattedMessage id='common.buttons.new' />
                </Button>
              </Link>
            </div>
          )}
        </>
      }
    >
      <List<Bookmark>
        resourceType="bookmark"
        loading={loading}
        error={error}
        meta={data?.meta}
        items={data?.items}
        changePageSize={changePageSize}
        paginate={paginate}
      />

      <FilterSelectionModal
        typeOfParent='bookmark'
        onApply={applyFilterAndSort}
        tags={tagsAsArray}
        sortableFields={['title', 'updatedAt']}
        sortBy={query.sortBy}
        sortDirection={query.sortDirection}
      />
    </Layout>
  );
}

export function BookmarksListPage() {
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
      <BookmarksList />
    </QueryContextProvider>
  )
}

