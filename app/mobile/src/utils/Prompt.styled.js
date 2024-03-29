import { StyleSheet } from 'react-native';
import { Colors } from 'constants/Colors';

export const styles = StyleSheet.create({
  modalOverlay: {
    width: '100%',
    height: '100%',
  },
  modalBase: {
    display: 'flex',
    width: '100%',
    height: '100%',
    alignItems: 'center',
    justifyContent: 'center',
    position: 'absolute',
    top: 0,
    left: 0,
  },
  modalContainer: {
    backgroundColor: Colors.modalBase,
    width: '80%',
    maxWidth: 400,
    display: 'flex',
    alignItems: 'center',
    borderRadius: 8,
    borderWidth: 1,
    borderColor: Colors.modalBorder,
  },
  modalHeader: {
    fontSize: 18,
    paddingTop: 16,
    color: Colors.labelText,
    fontFamily: 'Roboto',
  },
  modalMessage: {
    textAlign: 'center',
    fontSize: 14,
    color: Colors.descriptionText,
    paddingLeft: 16,
    paddingRight: 16,
    paddingBottom: 16,
  },
  closeButton: {
    marginTop: 8,
    marginBottom: 16,
    marginRight: 16,
    paddingTop: 8,
    paddingBottom: 8,
    paddingLeft: 32,
    paddingRight: 32,
    borderRadius: 4,
    borderWidth: 1,
    borderColor: Colors.cancelButton,
    backgroundColor: Colors.closeButton,
  },
  closeButtonText: {
    color: Colors.closeButtonText,
    fontFamily: 'Roboto',
  },
  cancelButton: {
    marginTop: 8,
    marginBottom: 16,
    marginRight: 16,
    paddingTop: 8,
    paddingBottom: 8,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: 4,
    width: 128,
    height: 32,
    backgroundColor: Colors.cancelButton,
  },
  cancelButtonText: {
    color: Colors.cancelButtonText,
    fontFamily: 'Roboto',
  },
  okButton: {
    marginTop: 8,
    marginBottom: 16,
    paddingTop: 8,
    paddingBottom: 8,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: 4,
    width: 128,
    height: 32,
    backgroundColor: Colors.primaryButton,
  },
  okButtonText: {
    color: Colors.primaryButtonText,
    fontFamily: 'Roboto',
  },
  modalButtons: {
    width: '100%',
    paddingLeft: 16,
    paddingRight: 16,
    paddingTop: 16,
    display: 'flex',
    alignItems: 'flex-end',
    justifyContent: 'flex-end',
    flexDirection: 'row',
  },
  centerModalButtons: {
    width: '100%',
    paddingLeft: 16,
    paddingRight: 16,
    paddingTop: 16,
    display: 'flex',
    alignItems: 'flex-end',
    justifyContent: 'center',
    flexDirection: 'row',
  },
});

