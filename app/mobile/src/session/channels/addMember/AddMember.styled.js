import { StyleSheet } from 'react-native';
import { Colors } from 'constants/Colors';

export const styles = StyleSheet.create({
  container: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    height: 48,
    alignItems: 'center',
    borderBottomWidth: 1,
    borderColor: Colors.itemDivider,
    paddingLeft: 8,
    paddingRight: 8,
  },
  detail: {
    paddingLeft: 12,
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    flexGrow: 1,
    flexShrink: 1,
  },
  space: {
    height: 64,
  },
  name: {
    color: Colors.text,
    fontSize: 14,
  },
  handle: {
    color: Colors.text,
    fontSize: 12,
  },
  connected: {
    width: 8,
    height: 8,
    borderRadius: 4,
    backgroundColor: Colors.connected,
  },
  requested: {
    width: 8,
    height: 8,
    borderRadius: 4,
    backgroundColor: Colors.requested,
  },
  connecting: {
    width: 8,
    height: 8,
    borderRadius: 4,
    backgroundColor: Colors.connecting,
  },
  pending: {
    width: 8,
    height: 8,
    borderRadius: 4,
    backgroundColor: Colors.pending,
  },
  confirmed: {
    width: 8,
    height: 8,
    borderRadius: 4,
    backgroundColor: Colors.confirmed,
  },
  track: {
    false: Colors.grey,
    true: Colors.background,
  },
  switch: {
    transform: [{ scaleX: .7 }, { scaleY: .7 }],
  },

})
