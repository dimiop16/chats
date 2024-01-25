import styled from 'styled-components';
import Colors from 'constants/Colors';

export const ChannelItemWrapper = styled.div`
  height: 48px;
  width: 100%;
  display: flex;
  align-items: center;
  border-bottom: 1px solid ${Colors.itemDivider};
  line-height: 16px;
  cursor: pointer;
  overflow: hidden;

  &:hover {
    background-color: ${Colors.formHover};
  }

  .active {
    background-color: ${Colors.profileForm};
    width: 100%;
    height: 100%;
    display: flex;
    align-item: center;
  }

  .idle {
    width: 100%;
  }

  .item {
    display: flex;
    flex-direction: row;
    align-items: center;
    min-width: 0;
    padding-left: 16px;
    padding-right: 16px;

    .avatar{
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 8px;
      font-size: 18px;
      flex-shrink: 0;
    }

    .logo {
      display: flex;
      align-items: center;
      justify-content: center;
      border: 1px solid ${Colors.grey};
      width: 32px;
      height: 32px;
      border-radius: 4px;
      font-size: 18px;
      flex-shrink: 0;
    }

    .details {
      flex-grow: 1;
      display: flex;
      flex-direction: column;
      padding-left: 16px;
      justify-content: center;
      min-width: 0;

      .subject {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .message {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        color: ${Colors.disabled};
      }
    }
  }
`

export const Markup = styled.div`
  position: absolute;
  right: 0;
  border-radius: 8px;
  background-color: ${Colors.background};
  width: 8px;
  height: 8px;
  margin-right: 8px;
`;
