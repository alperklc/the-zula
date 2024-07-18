import React from "react";
import LoadingIndicator from "../loadingIndicator";
import PageSizeSelector from "../pageSizeSelector";
import Pagination from "../pagination";
import { ListItem } from "../listItem";
import { PaginationMeta } from "../../types/Api";
import { Resource, ResourceType } from "../../types/resources";
import styles from "./index.module.css";

const PAGE_SIZE = 10

function getPaginationRange(count: number, page: number): string {
  if (count === 0) {
    return "";
  }

  const from = (page - 1) * PAGE_SIZE + 1;
  const to = Math.min(page * PAGE_SIZE, count);

  return `${from} - ${to} / ${count}`;
}

function ListBody<T>(
  { items, loading, error, render }:
    { items?: T[], loading: boolean, error: Error | string | null, render: (_: T) => React.ReactNode },
) {
  if (error) {
    return <p>An error occurred: {error?.toString()}</p>
  }

  if (loading) {
    return <LoadingIndicator />
  }

  if (!items || items?.length === 0) {
    return <span>Nothing found</span>
  }

  return <ul className={styles.container}>
    {items.map((item, index) => (
      <li key={index}>
        {render(item)}
      </li>
    ))}
  </ul>
}

function List<T>({ loading, error, meta, resourceType, items, changePageSize, paginate }: { loading: boolean, error: Error | string | null, meta?: PaginationMeta, resourceType: ResourceType, items?: T[], changePageSize: (_: number) => void, paginate: (_: number) => void }) {
  const maxPageNumber = Math.ceil((meta?.count || 0) / PAGE_SIZE)

  return (
    <div className={styles.container}>
      {!loading && <header className={styles.listHeader}>
        <span>
          <PageSizeSelector pageSize={meta?.pageSize ?? 0} onPageSizeSelected={changePageSize} />
        </span>
        {meta?.count ? (
          <span className={styles.count}>{getPaginationRange(meta?.count ?? 0, meta?.page ?? 1)}</span>
        ) : (
          <span />
        )}
      </header>}

      <div className={styles.body}>
        <ListBody<T>
          loading={loading}
          error={error}
          render={(item) => <ListItem resourceType={resourceType} item={item as Resource} />}
          items={items}
        />

        <footer className={styles.footer}>
          <span>
            {maxPageNumber > 0 ? (
              <Pagination
                onPageClicked={paginate}
                numberOfPages={maxPageNumber}
                currentPage={meta?.page ?? 1}
              />
            ) : (
              <span>&nbsp;</span>
            )}
          </span>
        </footer>
      </div>
    </div>)
}

export default List;
