import React, { useState, useEffect, useContext } from 'react';
import {render, act, screen, waitFor, fireEvent} from '@testing-library/react'
import { ProfileContextProvider } from 'context/ProfileContext';
import { ViewportContextProvider } from 'context/ViewportContext';
import { useListing } from 'session/listing/useListing.hook';
import * as fetchUtil from 'api/fetchUtil';

let listing = null;
function ListingView() {
  const { state, actions } = useListing();
  const [renderCount, setRenderCount] = useState(0);
  const [contacts, setContacts] = useState([]);

  listing = actions;
  useEffect(() => {

    const rendered = [];
    state.contacts.forEach(item => {
      rendered.push(
        <div key={item.guid} data-testid="contact">
          <span key={item.guid} data-testid={'contact-' + item.guid}>{ item.name }</span>
        </div>
      );
    });
    setContacts(rendered);
    setRenderCount(renderCount + 1);
  }, [state]);

  return (
    <div data-testid="contacts" count={renderCount}>
      { contacts }
    </div>
  );
}

function ListingTestApp() {
  return (
    <ProfileContextProvider>
      <ViewportContextProvider>
        <ListingView />
      </ViewportContextProvider>
    </ProfileContextProvider>
  );
}

let fetchListing;
const realFetchWithTimeout = fetchUtil.fetchWithTimeout;
const realFetchWithCustomTimeout = fetchUtil.fetchWithCustomTimeout;
beforeEach(() => {
  fetchListing = [];

  const mockFetch = jest.fn().mockImplementation((url, options) => {

    return Promise.resolve({
      json: () => Promise.resolve(fetchListing)
    });
  });
  fetchUtil.fetchWithTimeout = mockFetch;
  fetchUtil.fetchWithCustomTimeout = mockFetch;
});

afterEach(() => {
  fetchUtil.fetchWithTimeout = realFetchWithTimeout;
  fetchUtil.fetchWithCustomTimeout = realFetchWithCustomTimeout;
});

test('retrieve listing', async () => {
  render(<ListingTestApp />);

  await waitFor(() => {
    expect(listing).not.toBe(null);
  });

  fetchListing = [
    {
      guid: 'abc123',
      handle: 'tester',
      name: 'mr. tester',
      description: 'a tester',
      location: 'here',
      imageSet: false,
      version: '0.0.1',
      node: 'test.org',
    },
  ];

  await act(async () => {
    await listing.getListing();
  });

  await waitFor(async () => {
    expect(screen.getByTestId('contact-abc123').textContent).toBe('mr. tester');
  });

});



