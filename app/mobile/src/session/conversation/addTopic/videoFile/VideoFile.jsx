import { useRef, useEffect } from 'react';
import { TouchableOpacity, View } from 'react-native';
import Video from 'react-native-video';
import { useVideoFile } from './useVideoFile.hook';
import { styles } from './VideoFile.styled';
import MatIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import Colors from 'constants/Colors';

export function VideoFile({ path, setPosition, remove }) {

  const { state, actions } = useVideoFile();

  const video = useRef();

  useEffect(() => {
    if (video.current) {
      video.current.seek(state.position);
      setPosition(state.position);
    }
  }, [state.position]);

  const setInfo = ({ naturalSize, duration }) => {
    if (video.current) {
      video.current.seek(0);
    }
    actions.setInfo(naturalSize.width, naturalSize.height, duration);
  }

  return (
    <TouchableOpacity onPress={actions.setNextPosition} onLongPress={remove}>
      <Video source={{ uri: path }} style={{ width: 92 * state.ratio, height: 92, marginRight: 16 }} resizeMode={'cover'} paused={true}
        onLoad={setInfo} ref={(ref) => video.current = ref}
      />
      <View style={styles.overlay}>
        <MatIcons name="arrow-right" size={20} color={Colors.white} /> 
      </View>  
    </TouchableOpacity>
  );
}
