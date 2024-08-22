import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/react'
import Pagination from './index'

describe('Pagination', () => {
  [
    {
      numberOfPages: 8,
      currentPage: 3,
      expectedResult: ['1', '2', '3', '4', '5', '. . .', '8'],
    },
    {
      numberOfPages: 8,
      currentPage: 7,
      expectedResult: ['1', '. . .', '4', '5', '6', '7', '8'],
    },
    {
      numberOfPages: 14,
      currentPage: 9,
      expectedResult: ['1', '. . .', '8', '9', '10', '. . .', '14'],
    },
  ].map(({ numberOfPages, currentPage, expectedResult }) => {
    it('displays pagination buttons correctly', async () => {
      const pagination = render(
        <Pagination
          onPageClicked={() => ({})}
          numberOfPages={numberOfPages}
          currentPage={currentPage}
        />,
      )

      const childElements = await pagination.findAllByTestId('child-element')

      childElements.map((element, buttonIndex) => {
        expect(element.textContent).toEqual(expectedResult[buttonIndex])
      })
    })
  })
})
