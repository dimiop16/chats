import { useContext, useState, useEffect } from 'react';
import { useNavigate } from "react-router-dom";
import { getNodeStatus } from 'api/getNodeStatus';
import { setNodeStatus } from 'api/setNodeStatus';
import { getNodeConfig } from 'api/getNodeConfig';
import { AppContext } from 'context/AppContext';

export function useAdmin() {

  const [state, setState] = useState({
    password: '',
    placeholder: '',
    unclaimed: null,
    busy: false,
  });

  const navigate = useNavigate();
  const app = useContext(AppContext);

  const updateState = (value) => {
    setState((s) => ({ ...s, ...value }));
  }

  useEffect(() => {
    const check = async () => {
      try {
        const unclaimed = await getNodeStatus();
        updateState({ unclaimed });
      }
      catch(err) {
        console.log("getNodeStatus failed");
      }
    };
    check();
  }, []);

  const actions = {
    setPassword: (password) => {
      updateState({ password });
    },
    navUser: () => {
      navigate('/login');
    },
    login: async () => {
      if (!state.busy) {
        try {
          updateState({ busy: true });
          if (state.unclaimed === true) {
            await setNodeStatus(state.password);
          }
          await getNodeConfig(state.password);
          updateState({ busy: false });
          app.actions.setAdmin(state.password);          
        }
        catch(err) {
          console.log(err);
          updateState({ busy: false });
          throw new Error("login failed");
        }
      }
    },
  }

  return { state, actions };
}

