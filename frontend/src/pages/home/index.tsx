import { useIntl } from 'react-intl'
import Layout from '../../components/layout'
import MessageBox from '../../components/messageBox'
import PageContent from '../../components/pageContent'
import AllContentSection, { DashboardSection } from '../../components/dashboard/section'
import React from 'react'
import { useUI } from '../../contexts/uiContext'
import { useAuth } from '../../contexts/authContext'
import { Api, Insights } from '../../types/Api'
import { ActivityGraph } from '../../components/dashboard/activity-graph'

const emptyInsights: Insights = {
  numberOfNotes: 0,
  numberOfBookmarks: 0, 
  lastVisited: [],
  mostVisited: [],
  activityGraph: [],
}

const Dashboard = () => {
  const { formatMessage } = useIntl()
  const { isMobile } = useUI()

  const [loading, setLoading] = React.useState(true);
  const [data, setData] = React.useState<Insights>(emptyInsights);
  const [error, setError] = React.useState<string | null>(null);

  const { user } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}` } } })

  const fetch = async () => {
    try {
      setLoading(true);
      setError(null);

      const { data, error, status } = await api.api.getInsights(user?.profile.sub ?? "")

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
    fetch()
  }, [])

  const lastVisited = (data?.lastVisited || [])
  const mostVisited = (data?.mostVisited || [])

  const MostVisitedAndLastVisited = () => {
    return (
    <>
      {lastVisited?.length > 0 && (
        <DashboardSection
          title={formatMessage({ id: 'dashboard.titles.last_visited' })}
          rows={data?.lastVisited ?? []}
        />
      )}
      {mostVisited?.length > 0 && (
        <DashboardSection
          title={formatMessage({ id: 'dashboard.titles.most_visited' })}
          rows={data?.mostVisited ?? []}
        />
      )}
    </>
  )}

    return (
    <Layout narrow>
      {!loading && !data && error ? (
        <MessageBox type='error'>Error</MessageBox>
      ) : (
        <>
          <label>{formatMessage({ id: 'dashboard.titles.activity_graph' })}</label>

          <ActivityGraph data={data?.activityGraph ?? []} />

          {isMobile && (
            // on mobile always one column
            <PageContent loading={loading}>
              <>
                <MostVisitedAndLastVisited />
                <AllContentSection data={data} />
              </>
            </PageContent>
          )}

          {!isMobile && (
            // on desktop two columns BUT only one column is shown, if the most visited and last visited are empty
            <PageContent loading={loading}>
              <MostVisitedAndLastVisited />

              <AllContentSection data={data} />
            </PageContent>
          )}
        </>
      )}
    </Layout>
  )
}

export default Dashboard
