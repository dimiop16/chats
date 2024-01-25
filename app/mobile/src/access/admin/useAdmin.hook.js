import { useState, useEffect, useContext } from 'react';
import { useWindowDimensions } from 'react-native';
import { useNavigate } from 'react-router-dom';
import { AppContext } from 'context/AppContext';
import { getNodeStatus } from 'api/getNodeStatus';
import { setNodeStatus } from 'api/setNodeStatus';
import { getNodeConfig } from 'api/getNodeConfig';
import { getLanguageStrings } from 'constants/Strings';

export function useAdmin() {

  const navigate = useNavigate();
  const app = useContext(AppContext);

  const [state, setState] = useState({
    strings: getLanguageStrings(),
    busy: false,
    enabled: false,
    server: null,
    token: null,
    plainText: false,
    version: null,
    agree: false,
    showTerms: false,
  });

  const updateState = (value) => {
    setState((s) => ({ ...s, ...value }));
  }

  useEffect(() => {
    updateState({ version: app.state.version });
  }, [app]);

  const checkStatus = async () => {
    try {
      updateState({ unclaimed: status });
    }
    catch (err) {
      console.log("failed to check node status");
    }
  };

  useEffect(() => {
    if (state.token && state.server && !state.enabled) {
      updateState({ enabled: true });
    }
    if ((!state.token || !state.server) && state.enabled) {
      updateState({ enabled: false });
    }
  }, [state.server, state.token]);

  const actions = {
    setServer: (server) => {
      updateState({ server });
    },
    setToken: (token) => {
      updateState({ token });
    },
    login: () => {
      navigate('/login');
    },
    showPass: () => {
      updateState({ plainText: true });
    },
    hidePass: () => {
      updateState({ plainText: false });
    },
    showTerms: () => {
      updateState({ showTerms: true });
    },
    hideTerms: () => {
      updateState({ showTerms: false });
    },
    agree: (agree) => {
      updateState({ agree });
    },
    access: async () => {
      if (!state.busy) {
        try {
          updateState({ busy: true });
          const node = state.server.trim();
          const token = state.token;
          const unclaimed = await getNodeStatus(node);
          if (unclaimed) {
            await setNodeStatus(node, token);
          } 
          const config = await getNodeConfig(node, token);
          updateState({ server: node, busy: false });
          navigate('/dashboard', { state: { config, server: node, token }});
        }
        catch (err) {
          console.log(err);
          updateState({ busy: false });
          throw new Error("access failed");
        }
      }
    }
  };

  return { state, actions };
}

