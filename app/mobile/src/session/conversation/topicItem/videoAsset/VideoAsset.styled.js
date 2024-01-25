import { StyleSheet } from 'react-native';
import { Colors } from 'constants/Colors';

export const styles = StyleSheet.create({
  base: {
    position: 'relative',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    width: '100%',
    height: '100%',
  },
  container: {
    position: 'relative',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    width: '100%',
    height: '100%',
  },
  overlay: {
    position: 'absolute',
    height: '100%',
    width: '100%',
    backgroundColor: 'rgba(0, 0, 0, 0.4)',
  },
  control: {
    position: 'absolute',
    paddingRight: 8,
    paddingTop: 4,
  },
  thumb: {
    borderRadius: 4,
    opacity: 0.6,
    width: '100%',
    height: '100%',
  },
  main: {
    position: 'absolute',
    width: '100%',
    height: '100%',
  },
  close: {
    position: 'absolute',
    top: 0,
    right: 0,
    paddingRight: 16,
    paddingTop: 16,
  },
  loading: {
    position: 'absolute',
    display: 'flex',
    flexAlign: 'center',
    justifyContent: 'center',
    flexDirection: 'column',
  },
  decrypting: {
    fontVariant: ["tabular-nums"],
    paddingTop: 16,
    fontSize: 12,
    color: '#dddddd',
  },
})

