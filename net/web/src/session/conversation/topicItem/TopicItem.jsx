import { TopicItemWrapper } from './TopicItem.styled';
import { VideoAsset } from './videoAsset/VideoAsset';
import { AudioAsset } from './audioAsset/AudioAsset';
import { ImageAsset } from './imageAsset/ImageAsset';
import { BinaryAsset } from './binaryAsset/BinaryAsset';
import { Logo } from 'logo/Logo';
import { Space, Skeleton, Button, Modal, Input } from 'antd';
import { ExclamationCircleOutlined, DeleteOutlined, EditOutlined, FireOutlined, PictureOutlined } from '@ant-design/icons';
import { Carousel } from 'carousel/Carousel';
import { useTopicItem } from './useTopicItem.hook';

export function TopicItem({ host, contentKey, sealed, topic, update, remove }) {

  const [ modal, modalContext ] = Modal.useModal();
  const { state, actions } = useTopicItem(topic, contentKey);

  const removeTopic = () => {
    modal.confirm({
      title: 'Do you want to delete this message?',
      icon: <ExclamationCircleOutlined />,
      bodyStyle: { padding: 16 },
      okText: 'Yes, Delete',
      cancelText: 'No, Cancel',
      onOk:  async () => {
        try {
          await remove();
        }
        catch(err) {
          console.log(err);
          modal.error({
            title: 'Failed to Delete Message',
            content: 'Please try again.',
            bodyStyle: { padding: 16 },
          });
        }
      },
    });
  }

  const updateTopic = async () => {
    try {
      await update(state.message);
      actions.clearEditing();
    }
    catch(err) {
      console.log(err);
      modal.error({
        title: 'Failed to Update Message',
        content: 'Please try again.',
        bodyStyle: { padding: 16 },
      });
    }
  };

  const renderAsset = (asset, idx) => {
    if (asset.type === 'image') {
      return <ImageAsset asset={asset} />
    }
    if (asset.type === 'video') {
      return <VideoAsset asset={asset} />
    }
    if (asset.type === 'audio') {
      return <AudioAsset asset={asset} />
    }
    if (asset.type === 'binary') {
      return <BinaryAsset asset={asset} />
    }
    return <></>
  }

  return (
    <TopicItemWrapper>
      { modalContext }
      <div class="topic-header">
        <div class="avatar">
          <Logo width={32} height={32} radius={4} url={topic.imageUrl} />
        </div>
        <div class="info">
          <div class={ topic.nameSet ? 'set' : 'unset' }>{ topic.name }</div>
          <div>{ topic.createdStr }</div>
        </div>
        <div class="topic-options">
          <div class="buttons">
            { !sealed && topic.creator && (
              <div class="button edit" onClick={() => actions.setEditing(topic.text)}>
                <EditOutlined />
              </div>
            )}
            { (host || topic.creator) && (
              <div class="button remove" onClick={removeTopic}>
                <DeleteOutlined />
              </div>
            )}
          </div>
        </div>
      </div>
      { topic.status !== 'confirmed' && (
        <div class="skeleton">
          <Skeleton size={'small'} active={true} title={false} />
        </div>
      )}
      { topic.status === 'confirmed' && (
        <>
          { topic.assets?.length && (
            <>
              { topic.transform === 'error' && (
                <div class="asset-placeholder">
                  <FireOutlined style={{ fontSize: 32, color: '#ff8888' }} />
                </div>
              )}
              { topic.transform === 'incomplete' && (
                <div class="asset-placeholder">
                  <PictureOutlined style={{ fontSize: 32 }} />
                </div>
              )}
              { topic.transform === 'complete' && (
                <div class="topic-assets">
                  <Carousel pad={40} items={state.assets} itemRenderer={renderAsset} />
                </div>
              )}
            </>
          )}
          { sealed && (
            <div class="sealed-message">sealed message</div>
          )}
          { !sealed && !state.editing && (
            <div class="message">
              <div style={{ color: topic.textColor, fontSize: topic.textSize }}>{ topic.clickable }</div>
            </div>
          )}
          { state.editing && (
            <div class="editing">
              <Input.TextArea defaultValue={state.message} placeholder="message"
                style={{ resize: 'none', color: state.textColor, fontSize: state.textSize }}
                onChange={(e) => actions.setMessage(e.target.value)} rows={3} bordered={false}/>
              <div class="controls">
                <Space>
                  <Button onClick={actions.clearEditing}>Cancel</Button>
                  <Button type="primary" onClick={updateTopic}>Save</Button>
                </Space>
              </div>
            </div>
          )}
        </>
      )}
    </TopicItemWrapper>
  )
}

