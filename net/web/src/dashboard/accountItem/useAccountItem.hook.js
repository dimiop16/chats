import { useContext, useState, useEffect } from 'react';
import { getAccountImageUrl } from 'api/getAccountImageUrl';
import { setAccountStatus } from 'api/setAccountStatus';
import { addAccountAccess } from 'api/addAccountAccess';
import { ViewportContext } from 'context/ViewportContext';
import { AppContext } from 'context/AppContext';

export function useAccountItem(item, remove) {
  
  const [state, setState] = useState({
    statusBusy: false,
    removeBusy: false,
    accessBusy: false,
    showAccess: false,
  });
 
  const app = useContext(AppContext);
  const viewport = useContext(ViewportContext); 

  const updateState = (value) => {
    setState((s) => ({ ...s, ...value }));
  }

  useEffect(() => {

    updateState({
      disabled: item?.disabled,
      activeClass: item?.disabled ? 'inactive' : 'active',
      accountId: item?.accountId,
      name: item?.name,
      guid: item?.guid,
      handle: item?.handle,
      storage: Math.floor(item?.storageUsed > 1073741824 ? item?.storageUsed / 1073741824 : item?.storageUsed / 1048576),
      storageUnit: item?.storageUsed > 1073741824 ? "GB" : "MB",
      imageUrl: item?.imageSet ? getAccountImageUrl(app.state.adminToken, item?.accountId) : null,
    });
  }, [app.state.adminToken, item]); 

  useEffect(() => {
    updateState({ display: viewport.state.display });
  }, [viewport]);

  const actions = {
    setAccessLink: async () => {
      if (!state.accessBusy) {
        updateState({ accessBusy: true });
        try {
          const access = await addAccountAccess(app.state.adminToken, item.accountId);
          updateState({ accessToken: access, showAccess: true, accessBusy: false });
        }
        catch (err) {
          console.log(err);
          updateState({ accessBusy: false });
          throw new Error('failed to generate token');
        }
      }
    },
    setShowAccess: (showAccess) => {
      updateState({ showAccess });
    },
    remove: async () => {
      if (!state.removeBusy) {
        updateState({ removeBusy: true });
        try {
          await remove(state.accountId);
          updateState({ removeBusy: false });
        }
        catch(err) {
          console.log(err);
          updateState({ removeBusy: false });
          throw new Error('failed to remove account');
        }
      }
    },
    setStatus: async (disabled) => {
      if (!state.statusBusy) {
        updateState({ statusBusy: true });
        try {
          await setAccountStatus(app.state.adminToken, item.accountId, disabled);
          updateState({ statusBusy: false, disabled, activeClass: disabled ? 'inactive' : 'active' });
        }
        catch(err) {
          console.log(err);
          updateState({ statusBusy: false });
          throw new Error('failed to set account status'); 
        }
      }
    },
  };

  return { state, actions };
}
