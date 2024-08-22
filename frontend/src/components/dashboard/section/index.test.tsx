import { describe, vi, it, expect } from 'vitest';
import { render } from '@testing-library/react'
import { AllContentSection, DashboardSection } from './index'

vi.mock('react-intl', () => ({
  useIntl: () => ({ formatMessage: vi.fn() }),
}))

describe('DashboardSection', () => {
  it('WHEN an array of content is given THEN they are rendered correctly', () => {
    // arrange
    const rows = [
      {
          typename: 'Note',
          id: 'aabel2cGoibqDbkDIJf78',
          title: 'test note',
        },
        {  
          typename: 'Bookmark',
          id: 'AJ0K4sAPIqrG84AvDLOSL',
          name: 'kek.jpg',
        },
        {
          typename: 'Bookmark',
          id: '1FO5kHFBYZiFdxsdeblYu',
          title: 'Documentation - The Go Programming Language',
      },
    ]

    // act
    const { queryByTestId, queryAllByTestId } = render(
      <DashboardSection title='test' rows={rows} />,
    )
    const renderedContent = queryAllByTestId('content')

    // assert
    expect(queryByTestId('title')?.textContent).toBe('test')
    expect(renderedContent.length).toBe(rows.length)
    expect(renderedContent[0].textContent).toBe(rows[0].title)
    expect(renderedContent[1].textContent).toBe(rows[1].name)
    expect(renderedContent[2].textContent).toBe(rows[2].title)
  })
})

describe('AllContentSection', () => {
  it('WHEN all given numbers are 0 THEN the numbers are not displayed', () => {
    // act
    const { queryByTestId } = render(
      <AllContentSection data={{ numberOfBookmarks: 0, numberOfNotes: 0 }} />,
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
      <AllContentSection data={{ numberOfBookmarks: 5, numberOfNotes: 4 }} />,
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
