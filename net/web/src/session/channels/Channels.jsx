import { Modal, Input, List, Button, Tooltip } from 'antd';
import { ChannelsWrapper } from './Channels.styled';
import { CommentOutlined, SearchOutlined } from '@ant-design/icons';
import { useChannels } from './useChannels.hook';
import { ChannelItem } from './channelItem/ChannelItem';
import { AddChannel } from './addChannel/AddChannel';

export function Channels({ open, active }) {

  const { state, actions } = useChannels();

  const added = (id) => {
    actions.clearShowAdd();
    open(id);
  };

  return (
    <ChannelsWrapper>
      <div class="search">
        <div class="filter">
          <Input bordered={false} allowClear={true} placeholder="Topics" prefix={<SearchOutlined />}
              spellCheck="false" onChange={(e) => actions.onFilter(e.target.value)} />
        </div>
        { state.display === 'small' && (
          <div class="inline">
            <Button type="primary" disabled={!state.allowAdd} icon={<CommentOutlined />} onClick={actions.setShowAdd}>New</Button>
          </div>
        )}
      </div>
      <div class="view">
        { state.channels.length > 0 && (
          <List local={{ emptyText: '' }} itemLayout="horizontal" dataSource={state.channels} gutter="0"
            renderItem={item => (
              <ChannelItem item={item} openChannel={open}
                  active={active.card === item.cardId && active.channel === item.channelId} />
            )}
          />
        )}
        { state.channels.length === 0 && (
          <div class="empty">No Topics</div>
        )}
      </div>
      { state.display !== 'small' && (
        <div class="bar">
          <Tooltip placement="right" title={ state.allowAdd ? '' : 'Account Sealing Key Required' }>
            <Button type="primary" disabled={!state.allowAdd} icon={<CommentOutlined />} onClick={actions.setShowAdd}>New Topic</Button>
          </Tooltip>
        </div>
      )}
      <Modal bodyStyle={{ padding: 16 }} title="New Topic" centered visible={state.showAdd && state.allowAdd} footer={null} destroyOnClose={true}
          onCancel={actions.clearShowAdd}>
        <AddChannel added={added} cancelled={actions.clearShowAdd} />
      </Modal>
    </ChannelsWrapper>
  );
}

