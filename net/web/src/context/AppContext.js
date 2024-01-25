import { createContext } from 'react';
import { useAppContext } from './useAppContext.hook';

export const AppContext = createContext({});

export function AppContextProvider({ children }) {
  const { state, actions } = useAppContext();
  return (
    <AppContext.Provider value={{ state, actions }}>
      {children}
    </AppContext.Provider>
  );
}

