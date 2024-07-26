import { render } from '@testing-library/react'
import { AllContentSection, DashboardSection, Content, DashboardContentRow } from './index'

jest.mock('react-intl', () => ({
  useIntl: () => ({ formatMessage: jest.fn() }),
}))

describe('DashboardSection', () => {
  it('WHEN an array of content is given THEN they are rendered correctly', () => {
    // arrange
    const rows = [
      {
        content: {
          __typename: 'Note',
          uid: 'aabel2cGoibqDbkDIJf78',
          title: 'test note',
        },
      },
      {
        content: {
          __typename: 'File',
          uid: 'AJ0K4sAPIqrG84AvDLOSL',
          name: 'kek.jpg',
        },
      },
      {
        content: {
          __typename: 'Bookmark',
          uid: '1FO5kHFBYZiFdxsdeblYu',
          title: 'Documentation - The Go Programming Language',
        },
      },
    ] as DashboardContentRow[]

    // act
    const { queryByTestId, queryAllByTestId } = render(
      <DashboardSection title='test' rows={rows} />,
    )
    const renderedContent = queryAllByTestId('content')

    // assert
    expect(queryByTestId('title')?.textContent).toBe('test')
    expect(renderedContent.length).toBe(rows.length)
    expect(renderedContent[0].textContent).toBe(rows[0].content.title)
    expect(renderedContent[1].textContent).toBe(rows[1].content.name)
    expect(renderedContent[2].textContent).toBe(rows[2].content.title)
  })
})

describe('AllContentSection', () => {
  it('WHEN all given numbers are 0 THEN the numbers are not displayed', () => {
    // act
    const { queryByTestId } = render(
      <AllContentSection data={{ numberOfBookmarks: 0, numberOfFiles: 0, numberOfNotes: 0 }} />,
    )
    const [numberOfNotes, numberOfFiles, numberOfBookmarks] = [
      'numberOfNotes',
      'numberOfFiles',
      'numberOfBookmarks',
    ].map((testId) => queryByTestId(testId))

    // assert
    expect(numberOfNotes?.textContent).toBe('')
    expect(numberOfFiles?.textContent).toBe('')
    expect(numberOfBookmarks?.textContent).toBe('')
  })

  it('WHEN some of given numbers are NOT 0 THEN they are displayed', () => {
    // act
    const { queryByTestId } = render(
      <AllContentSection data={{ numberOfBookmarks: 5, numberOfFiles: 0, numberOfNotes: 4 }} />,
    )
    const [numberOfNotes, numberOfFiles, numberOfBookmarks] = [
      'numberOfNotes',
      'numberOfFiles',
      'numberOfBookmarks',
    ].map((testId) => queryByTestId(testId))

    // assert
    expect(numberOfNotes?.textContent).toBe(' (4)')
    expect(numberOfFiles?.textContent).toBe('')
    expect(numberOfBookmarks?.textContent).toBe(' (5)')
  })
})
