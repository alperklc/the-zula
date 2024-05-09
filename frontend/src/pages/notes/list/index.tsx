import React from "react";
import { Link, useSearchParams } from "react-router-dom";
import { FormattedMessage } from "react-intl";
import Layout, { styles as layoutStyles } from "../../../components/layout";
import List from "../../../components/list";
import SearchInput from "../../../components/search";
import { useAuth } from "../../../contexts/authContext";
import { Api, Note, PaginationMeta } from "../../../types/Api";
import { useUI } from "../../../contexts/uiContext";
import Breadcrumbs from "../../../components/breadcrumbs";
import Button from "../../../components/form/button";
import { Query, QueryContext, QueryContextProvider } from "../../../contexts/queryContext";
import { filterEmptyValues } from "../../../utils/filter"
import MobileHeader from "../../../components/mobileHeader";
import useModal from "../../../components/modal";
import FilterSelector, { FilterSelectorModalProps } from "../../../components/filterSelectorModal";
import Icons from "../../../components/icons";

function NotesList() {
  const { isMobile } = useUI()
  const [FilterSelectionModal, openFilterSelection] =
    useModal<FilterSelectorModalProps>(FilterSelector)

  const { query, applyFilterAndSort, tagsAsArray, changePageSize, updateQuery, paginate } = React.useContext(QueryContext)
  const [loading, setLoading] = React.useState(true);
  const [data, setData] = React.useState<{ meta?: PaginationMeta, items?: Note[] }>();
  const [error, setError] = React.useState<string | null>(null);
  
  const { user } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}` } } })

  const fetchAndUpdate = async (query: Query) => {
    try {
      setLoading(true);
      setError(null);

      const filteredQuery = filterEmptyValues(query)
      const { data, status } = await api.api.v1NotesList(filteredQuery)

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
  }, [query.page, query.pageSize, query.q, query.sortBy, query.sortDirection]);

  return (
    <Layout
      customHeader={
        isMobile
          ? (props: unknown) => <MobileHeader onFilterSelectionOpen={openFilterSelection} {...props}/>
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
                  <Button muted onClick={openFilterSelection} className={layoutStyles.filterButton}>
                    <Icons.Filter width={18} height={18} />
                    {tagsAsArray.length > 0 && (
                      <span className={layoutStyles.badgeOnButton}>{tagsAsArray.length}</span>
                    )}
                  </Button>
                </div>
              </div>

              <Link to='/notes/create'>
                <Button>
                  <FormattedMessage id='common.buttons.new_note' />
                </Button>
              </Link>
            </div>
          )}
        </>
      }
    >
      <List<Note>
        resourceType="note"
        loading={loading}
        error={error}
        meta={data?.meta}
        items={data?.items}
        changePageSize={changePageSize}
        paginate={paginate}
      />

      <FilterSelectionModal
        typeOfParent='note'
        onApply={applyFilterAndSort}
        tags={tagsAsArray}
        sortableFields={['title', 'updatedAt']}
        sortBy={query.sortBy}
        sortDirection={query.sortDirection}
      />
    </Layout>
  );
}

export function ListPage() {
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
      <NotesList />
    </QueryContextProvider>
  )
}

